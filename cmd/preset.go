package cmd

import (
	"fmt"
	"strings"

	"github.com/jheddings/ccglow/internal/preset"
	"github.com/spf13/cobra"
)

var presetCmd = &cobra.Command{
	Use:   "preset",
	Short: "Inspect available presets",
}

var presetListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available preset names",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		names := preset.List()
		if len(names) == 0 {
			return fmt.Errorf("no presets found")
		}
		fmt.Fprintln(cmd.OutOrStdout(), strings.Join(names, "\n"))
		return nil
	},
}

var presetShowCmd = &cobra.Command{
	Use:   "show <name>",
	Short: "Show the JSON config for a preset",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := preset.Dump(args[0])
		if err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), string(data))
		return nil
	},
}

func init() {
	presetCmd.AddCommand(presetListCmd, presetShowCmd)
	rootCmd.AddCommand(presetCmd)
}
