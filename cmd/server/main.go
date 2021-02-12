package main

import (
	"accounts/internal/server/grpc"
	"accounts/internal/server/repository"
	"accounts/internal/server/repository/inmem"
	_pg "accounts/internal/server/repository/pg"
	"accounts/internal/server/service"
	"accounts/pkg/cache/lru"
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/go-pg/pg/v10"
)

var (
	addr            string
	cacheSize       int
	pgURL           string
	pollingInterval string
)

func main() {
	flag.StringVar(&addr, "addr", ":8080", "listening address")
	flag.IntVar(&cacheSize, "cache-size", 10, "cache size")
	flag.StringVar(&pgURL, "pg-url", os.Getenv("PG_URL"), "the postgresql connection string. If omitted, the PG_URL"+
		"environment variable is searched for. If PG_URL not specified, then we use the in-memory storage")
	flag.StringVar(&pollingInterval, "polling-interval", os.Getenv("POLLING_INTERVAL"), "polling interval by the"+
		" statistics collection service. If omitted, the POLLING_INTERVAL environment variable is searched for."+
		" If POLLING_INTERVAL not specified, the default value is '30s'")
	flag.Parse()

	var repo repository.Accounts
	if pgURL == "" {
		repo = inmem.NewAccountsRepo()
	} else {
		opt, err := pg.ParseURL(pgURL)
		if err != nil {
			panic(err)
		}

		db := pg.Connect(opt)
		defer db.Close()

		//проверим подключение
		if _, err := db.Exec("SELECT 1"); err != nil {
			panic(err)
		}

		repo = _pg.NewAccountsRepo(db)
	}

	var _pollingInterval time.Duration
	if pollingInterval == "" {
		_pollingInterval = 30 * time.Second
	} else {
		var err error
		if _pollingInterval, err = time.ParseDuration(pollingInterval); err != nil {
			panic(err)
		}
	}

	statisticsSvc := service.NewStatisticsSvc(context.Background(), _pollingInterval)

	cache := lru.NewCache(cacheSize)
	accountsSvc := service.NewAccountsSvc(repo, cache)
	accountsSrv := grpc.
		NewServer(addr, accountsSvc, statisticsSvc).
		WithUnaryInterceptors(
			func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
				if strings.Contains(info.FullMethod, "GetAmount") {
					statisticsSvc.IncReadOperations()
				}
				if strings.Contains(info.FullMethod, "AddAmount") {
					statisticsSvc.IncWriteOperations()
				}

				return handler(ctx, req)
			},
		)

	doneCh := make(chan os.Signal, 1)
	signal.Notify(doneCh, os.Interrupt)

	errCh := make(chan error)

	go func() {
		if err := accountsSrv.Start(); err != nil {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		log.Println(err)
	case <-doneCh:
	}

	accountsSrv.GracefulStop()
}
