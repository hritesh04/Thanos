package cmd

import (
	"log/slog"
	"os"

	"github.com/hritesh04/thanos/core"
	"github.com/hritesh04/thanos/pkg/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	Lb      *core.LoadBalancer
	cfgFile string
	cfg     config.Config
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start Thanos Server for Load Balancing",
	Long:  "Start Thanos Server for Load Balancing",
	Run: func(cmd *cobra.Command, args []string) {
		initConfig()
		Lb = core.NewLoadBalancer(cfg)
		Lb.AddServer(cfg.Servers[0])
		Lb.Start()
	},
}

func init() {
	startCmd.Flags().StringVarP(&cfgFile, "config", "c", "", "config file default is example.yml")
	// startCmd.MarkFlagRequired("config")
	rootCmd.AddCommand(startCmd)
}

func initConfig() {
	if cfgFile != "" {
		config, err := os.ReadFile(cfgFile)
		if err != nil {
			slog.Error("config file cant be parsed")
		}
		if err := yaml.Unmarshal(config, &cfg); err != nil {
			slog.Error("config file cant be parsed")
		}
	} else {
		slog.Info("Config file not provided using default config")
		config, err := os.ReadFile("./example.yml")
		if err != nil {
			slog.Error("config file cant be parsed")
		}
		if err := yaml.Unmarshal(config, &cfg); err != nil {
			slog.Error("config file cant be parsed")
		}
	}
}
