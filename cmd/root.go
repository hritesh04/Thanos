package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "Thanos",
	Short: "Thanos is a very reliable Load Balancer",
	Long: `A Fast and Powerful Load Balancer to take of
		all the traffic and distribute them between nodes.
		
		"When I'm done, there will be no overloading on a server.
		Perfect balance, as all things should be...." ~ Thanos
		`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		fmt.Println("Hello from root")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
