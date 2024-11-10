package cmd

import (
	"fmt"
	"log"

	"github.com/elliot14A/jondev/domain/pkg"
	"github.com/elliot14A/jondev/interfaces/grpc"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "run jondev grpc server",
	Run: func(cmd *cobra.Command, args []string) {
		if err := serve(); err != nil {
			log.Fatalf("❌ Failed to run grpc server: %v", err)
		}
	},
}

func serve() error {
	config, err := pkg.LoadConfig()
	if err != nil {
		return fmt.Errorf("❌ Error loading config: %v", err)
	}
	return grpc.RunGrpcServer(*config)
}
