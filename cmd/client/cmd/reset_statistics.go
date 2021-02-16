package cmd

import (
	"accounts/internal/client"
	"accounts/pkg/log"
	"context"

	"github.com/spf13/cobra"
)

var resetStatCmd = &cobra.Command{
	Use:   "reset-stat",
	Short: "Resetting the statistics of transactions on the accounts",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := readConfig()
		if err != nil {
			log.Error(err.Error())
			return
		}

		client := client.NewStatisticsServiceClient(cfg.Client.Addr)
		if err = client.ResetStatistics(context.Background()); err != nil {
			log.Error(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(resetStatCmd)
}
