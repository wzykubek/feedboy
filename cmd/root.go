package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go.wzykubek.xyz/feedboy/pkg/server"
)

var (
	port       string
	schemesDir string
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&port, "port", "p", "8080", "listening port")
	rootCmd.PersistentFlags().StringVarP(&schemesDir, "directory", "d", "./schemes", "directory containing schemes")

	rootCmd.Flags().SortFlags = false
	rootCmd.PersistentFlags().SortFlags = false
}

var rootCmd = &cobra.Command{
	Use:  "feedboy",
	Long: "Self-hosted, template based RSS feed generator for websites without one",
	Run: func(cmd *cobra.Command, args []string) {
		srv := server.Server{
			Port:      port,
			SchemeDir: schemesDir,
		}

		srv.LoadSchemes()

		srv.Start()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
