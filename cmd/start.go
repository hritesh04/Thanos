package cmd

import (
	"net/http"
	"os"

	"github.com/hritesh04/thanos/internal"
	"github.com/hritesh04/thanos/internal/proxy"
	"github.com/hritesh04/thanos/internal/types"
	"github.com/hritesh04/thanos/pkg/config"
	"github.com/hritesh04/thanos/pkg/logger"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	Lb      types.IBalancer
	cfgFile string
	cfg     config.Config
	LogFile *os.File
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start Thanos Server for Load Balancing",
	Long:  "Start Thanos Server for Load Balancing",
	Run: func(cmd *cobra.Command, args []string) {
		LogFile = logger.InitLogger()
		defer LogFile.Close()
		logger.Log.Info("Parsing config file")
		initConfig()
		logger.Log.Info("Config file parsed successfully")
		logger.Log.Info("Creating Proxies")
		Lb = internal.NewLoadBalancer(cfg, proxy.NewReverseProxy)
		logger.Log.Info("Proxies created successfully")
		logger.Log.Info("Starting Thanos Server")
		http.HandleFunc("/", Lb.Serve)
		logger.Log.Info("Server Running at port " + cfg.ListenPort)
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
			logger.Log.Error("config file cant be parsed")
		}
		if err := yaml.Unmarshal(config, &cfg); err != nil {
			logger.Log.Error("config file cant be parsed")
		}
	} else {
		logger.Log.Info("Config file not provided using default config")
		config, err := os.ReadFile("./example.yml")
		if err != nil {
			logger.Log.Error("config file cant be parsed")
		}
		if err := yaml.Unmarshal(config, &cfg); err != nil {
			logger.Log.Error("config file cant be parsed")
		}
	}
}
