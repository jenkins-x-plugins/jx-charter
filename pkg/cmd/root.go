package cmd

import (
	"github.com/jenkins-x-plugins/jx-charter/pkg/cmd/run"
	"github.com/jenkins-x-plugins/jx-charter/pkg/cmd/update"
	"github.com/jenkins-x-plugins/jx-charter/pkg/cmd/version"
	"github.com/jenkins-x-plugins/jx-charter/pkg/rootcmd"
	"github.com/jenkins-x/jx-helpers/v3/pkg/cobras"
	"github.com/jenkins-x/jx-logging/v3/pkg/log"
	"github.com/spf13/cobra"
)

// Main creates the new command
func Main() *cobra.Command {
	cmd := &cobra.Command{
		Annotations: map[string]string{
			cobra.CommandDisplayNameAnnotation: rootcmd.TopLevelCommand,
		},
		Short: "commands for generating Helm Chart CRDs for better reporting and insight",
		Run: func(cmd *cobra.Command, args []string) {
			err := cmd.Help()
			if err != nil {
				log.Logger().Errorf(err.Error())
			}
		},
	}
	cmd.AddCommand(cobras.SplitCommand(run.NewCmdRun()))
	cmd.AddCommand(cobras.SplitCommand(update.NewCmdUpdate()))
	cmd.AddCommand(cobras.SplitCommand(version.NewCmdVersion()))
	return cmd
}
