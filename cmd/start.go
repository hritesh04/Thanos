package cmd

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/hritesh04/thanos/internal"
	"github.com/hritesh04/thanos/internal/proxy"
	"github.com/hritesh04/thanos/internal/types"
	"github.com/hritesh04/thanos/pkg/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	Lb      types.IBalancer
	cfgFile string
	cfg     config.Config
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start Thanos Server for Load Balancing",
	Long:  "Start Thanos Server for Load Balancing",
	Run: func(cmd *cobra.Command, args []string) {
		initConfig()
		Lb = internal.NewLoadBalancer(cfg, proxy.NewReverseProxy)
		http.HandleFunc("/", Lb.Serve)
		http.ListenAndServe(":"+cfg.ListenPort, nil)
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
