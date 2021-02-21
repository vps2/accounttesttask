package cmd

import (
	"os"
	"path/filepath"

	"github.com/vps2/accounttesttask/internal/client/config"
	"github.com/vps2/accounttesttask/pkg/log"

	"github.com/spf13/cobra"
	"gopkg.in/natefinch/lumberjack.v2"
)

var rootCmd = &cobra.Command{
	Use: "client",
}

func Execute() {
	initLog()

	rootCmd.Execute()
}

func init() {
	executableDir := filepath.Dir(os.Args[0])
	defaultCfgFile := filepath.Join(executableDir, "config/config.yml")
	rootCmd.PersistentFlags().String("cfg-file", defaultCfgFile, "path to config file")

	rootCmd.PersistentFlags().String("log-file", "", "path to a log file")
}

func readConfig() (*config.Config, error) {
	cfgFile, _ := rootCmd.PersistentFlags().GetString("cfg-file")
	cfg, err := config.New(cfgFile)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func initLog() {
	logFile, _ := rootCmd.PersistentFlags().GetString("log-file")
	if logFile != "" {
		log.SetOutput(&lumberjack.Logger{
			Filename:   logFile,
			MaxSize:    1,
			MaxBackups: 5,
			LocalTime:  true,
		})
	}
}
