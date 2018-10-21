package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "openstack-cli",
	Short: "OpenStack CLI",
	Long:  `OpenStack CLI...`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

// Execute root cmd
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
