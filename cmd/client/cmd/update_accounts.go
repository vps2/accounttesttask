package cmd

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/vps2/accounttesttask/internal/client"
	"github.com/vps2/accounttesttask/pkg/log"

	"github.com/spf13/cobra"
)

var updateAccountsCmd = &cobra.Command{
	Use:   "update-accounts",
	Short: "Updating accounts",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := readConfig()
		if err != nil {
			log.Error(err.Error())
			return
		}

		ctx, cancel := context.WithCancel(context.Background())

		var wg sync.WaitGroup

		for i := 0; i < cfg.Client.Readers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				client := client.NewAccountServiceClient(cfg.Client.Addr, cfg.Client.Keys, client.OpRead)
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

				client := client.NewAccountServiceClient(cfg.Client.Addr, cfg.Client.Keys, client.OpWrite)
				if err := client.Run(ctx); err != nil {
					log.Error(err.Error())
					return
				}
			}()
		}

		doneCh := make(chan os.Signal, 1)
		signal.Notify(doneCh, os.Interrupt)

		idle, _ := cmd.Flags().GetDuration("idle")

		//ожидаем нажатия Ctrl+C или наступления таймаута
		select {
		case <-time.After(idle):
		case <-doneCh:
		}

		cancel()

		wg.Wait()

		log.Info("done")
	},
}

func init() {
	updateAccountsCmd.Flags().Duration("idle", 15*time.Second, "the waiting time (in seconds) before program exit")

	rootCmd.AddCommand(updateAccountsCmd)
}
