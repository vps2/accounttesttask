package main

import (
	"accounts/internal/client"
	"accounts/internal/client/config"
	"accounts/pkg/log"
	"context"
	"flag"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	executableDir := filepath.Dir(os.Args[0])
	defaultCfgFile := filepath.Join(executableDir, "config/config.yml")

	idle := flag.Duration("idle", 15*time.Second, "the waiting time (in seconds) before program exit")
	cfgFile := flag.String("cfg-file", defaultCfgFile, "path to config file")
	logFile := flag.String("log-file", "", "path to a log file")

	flag.Parse()

	if *logFile != "" {
		log.SetOutput(&lumberjack.Logger{
			Filename:   *logFile,
			MaxSize:    1,
			MaxBackups: 5,
			LocalTime:  true,
		})
	}

	cfg, err := config.New(*cfgFile)
	if err != nil {
		log.Fatal(err.Error())
	}

	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup

	for i := 0; i < cfg.Client.Readers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			client := client.New(cfg.Client.Addr, cfg.Client.Keys, client.OpRead)
			if err := client.Run(ctx); err != nil {
				log.Error(err.Error())
				return
			}
		}()
	}
	for i := 0; i < cfg.Client.Writers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			client := client.New(cfg.Client.Addr, cfg.Client.Keys, client.OpWrite)
			if err := client.Run(ctx); err != nil {
				log.Error(err.Error())
				return
			}
		}()
	}

	doneCh := make(chan os.Signal, 1)
	signal.Notify(doneCh, os.Interrupt)

	//ожидаем нажатия Ctrl+C или наступления таймаута
	select {
	case <-time.After(*idle):
	case <-doneCh:
	}

	cancel()

	wg.Wait()
	log.Info("done")
}
