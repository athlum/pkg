package utils

import (
	"github.com/spf13/cobra"
	"strings"
)

func CMDIsHelp(cmd *cobra.Command, args []string) bool {
	return (len(args) == 1 && args[0] == "help") || strings.HasPrefix(cmd.Use, "help")
}

func RunningHelp(cmd *cobra.Command, args []string) bool {
	val := CMDIsHelp(cmd, args)
	if val {
		cmd.UsageFunc()(cmd)
	}
	return val
}
