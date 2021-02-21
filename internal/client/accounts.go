package client

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/vps2/accounttesttask/internal/api"
	"github.com/vps2/accounttesttask/pkg/log"

	"google.golang.org/grpc"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Operation int

const (
	OpRead Operation = iota + 1
	OpWrite
)

func (t Operation) String() string {
	switch t {
	case OpRead:
		return "Reader"
	case OpWrite:
		return "Writer"
	}

	return "Unknown"
}

var id int32

type AccountsServiceClient struct {
	id        int32
	addr      string
	keys      []int
	operation Operation
	//
	trigger *sync.WaitGroup
}

func NewAccountServiceClient(addr string, keys []int, op Operation) *AccountsServiceClient {
	return &AccountsServiceClient{
		id:        atomic.AddInt32(&id, 1),
		addr:      addr,
		keys:      keys,
		operation: op,
	}
}

func (c *AccountsServiceClient) WithTrigger(trigger *sync.WaitGroup) *AccountsServiceClient {
	c.trigger = trigger

	return c
}

func (c *AccountsServiceClient) Run(ctx context.Context) error {
	conn, err := grpc.Dial(c.addr, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("[%d] %w", c.id, err)
	}
	defer conn.Close()

	client := api.NewAccountsServiceClient(conn)

	if c.trigger != nil {
		c.trigger.Wait()
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			err := c.doJob(ctx, client)
			if err != nil {
				log.Errorf("[%d]\t%s\n", c.id, err)
			}
		}
	}
}

func (c *AccountsServiceClient) doJob(ctx context.Context, client api.AccountsServiceClient) error {
	idx := rand.Int31n(int32(len(c.keys)))
	balanceId := c.keys[idx]

	switch c.operation {
	case OpRead:
		resp, err := client.GetAmount(ctx, &api.GetRequest{BalanceId: int32(balanceId)})
		if err == nil {
			log.Infof("[%d]\taccount_%d\trequested balance: %d\n", c.id, balanceId, resp.Amount)
		}

		return err
	case OpWrite:
		var minBound int64 = -10
		var maxBound int64 = 11

		amount := rand.Int63n(maxBound-minBound) + minBound

		_, err := client.AddAmount(ctx, &api.AddRequest{BalanceId: int32(balanceId), Value: amount})
		log.Infof("[%d]\taccount_%d\tadd amount %d", c.id, balanceId, amount)

		return err
	}

	return errors.New("unknown operation")
}
