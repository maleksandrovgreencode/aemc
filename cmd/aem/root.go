package main

import (
	"github.com/spf13/cobra"
	"github.com/wttech/aemc/pkg/cfg"
	"github.com/wttech/aemc/pkg/instance"
	"strings"
)

func (c *CLI) rootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "aem",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			c.onStart()
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			c.onEnd()
			return nil
		},
	}
	cmd.AddCommand(c.versionCmd())
	cmd.AddCommand(c.initCmd())
	cmd.AddCommand(c.configCmd())
	cmd.AddCommand(c.instanceCmd())
	cmd.AddCommand(c.osgiCmd())
	cmd.AddCommand(c.pkgCmd())
	cmd.AddCommand(c.repoCmd())
	cmd.AddCommand(c.replCmd())
	cmd.AddCommand(c.cryptoCmd())
	cmd.AddCommand(c.fileCmd())
	c.rootFlags(cmd)
	return cmd
}

func (c *CLI) rootFlags(cmd *cobra.Command) {
	cv := c.config.Values()

	cmd.PersistentFlags().String("input-format", cv.GetString("input.format"), "Controls input format ("+strings.Join(cfg.InputFormats(), "|")+")")
	_ = cv.BindPFlag("input.format", cmd.PersistentFlags().Lookup("input-format"))

	cmd.PersistentFlags().String("input-file", cv.GetString("input.file"), "Provides input as file path")
	_ = cv.BindPFlag("input.file", cmd.PersistentFlags().Lookup("input-file"))

	cmd.PersistentFlags().String("input-string", cv.GetString("input.string"), "Provides input as string")
	_ = cv.BindPFlag("input.string", cmd.PersistentFlags().Lookup("input-string"))

	cmd.PersistentFlags().String("output-value", cv.GetString("output.value"),
		"Limits output to single variable")
	_ = cv.BindPFlag("output.value", cmd.PersistentFlags().Lookup("output-value"))

	cmd.PersistentFlags().String("output-format", cv.GetString("output.format"), "Controls output format ("+strings.Join(cfg.OutputFormats(), "|")+")")
	_ = cv.BindPFlag("output.format", cmd.PersistentFlags().Lookup("output-format"))

	cmd.PersistentFlags().String("output-log-file", cv.GetString("output.log.file"), "Controls output file path")
	_ = cv.BindPFlag("output.log.file", cmd.PersistentFlags().Lookup("output-log-file"))

	cmd.PersistentFlags().String("output-log-mode", cv.GetString("output.log.mode"), "Controls where outputs and logs should be written to when format is \"text\""+(strings.Join(cfg.OutputLogModes(), "|")+")"))
	_ = cv.BindPFlag("output.log.mode", cmd.PersistentFlags().Lookup("output-log-mode"))

	cmd.PersistentFlags().StringP("instance-url", "U", cv.GetString("instance.adhoc_url"), "Use only AEM instance at ad-hoc specified URL")
	_ = cv.BindPFlag("instance.adhoc_url", cmd.PersistentFlags().Lookup("instance-url"))

	cmd.PersistentFlags().StringP("instance-id", "I", cv.GetString("instance.filter.id"), "Use only AEM instance configured with the exact ID")
	_ = cv.BindPFlag("instance.filter.id", cmd.PersistentFlags().Lookup("instance-id"))

	cmd.PersistentFlags().BoolP("instance-author", "A", cv.GetBool("instance.filter.authors"), "Use only AEM author instance")
	_ = cv.BindPFlag("instance.filter.authors", cmd.PersistentFlags().Lookup("instance-author"))

	cmd.PersistentFlags().BoolP("instance-publish", "P", cv.GetBool("instance.filter.publishes"), "Use only AEM publish instance")
	_ = cv.BindPFlag("instance.filter.publishes", cmd.PersistentFlags().Lookup("instance-publish"))

	cmd.MarkFlagsMutuallyExclusive("instance-author", "instance-publish")

	cmd.PersistentFlags().String("instance-processing", cv.GetString("instance.processing_mode"), "Controls processing mode for instances ("+(strings.Join(instance.ProcessingModes(), "|")+")"))
	_ = cv.BindPFlag("instance.processing_mode", cmd.PersistentFlags().Lookup("instance-processing"))
}
