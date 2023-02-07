package cmd

import "github.com/spf13/cobra"

func New() *cobra.Command {
	root := &cobra.Command{
		Use:   "aauth <subcommand>",
		Short: "aauth Daemon",
		Long:  `aauth is a gRPC microservice. More info at https://github.com/web-programming-fall-2022/airline-auth`,
		Run:   nil,
	}
	addServeCmd(root)
	return root
}
