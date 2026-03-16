package cmd

import (
	"fmt"
	"io"
	"os"
	"runtime/debug"

	"github.com/jheddings/ccglow/internal/config"
	"github.com/jheddings/ccglow/internal/preset"
	"github.com/jheddings/ccglow/internal/provider"
	"github.com/jheddings/ccglow/internal/render"
	"github.com/jheddings/ccglow/internal/segment"
	"github.com/jheddings/ccglow/internal/session"
	"github.com/jheddings/ccglow/internal/style"
	"github.com/jheddings/ccglow/internal/types"
	"github.com/spf13/cobra"
)

var (
	presetName string
	configPath string
	format     string
	tee        string
)

var rootCmd = &cobra.Command{
	Use:   "ccglow",
	Short: "Composable statusline for Claude Code",
	Long:  "Reads session JSON from stdin, outputs styled statusline to stdout.",
	RunE: func(cmd *cobra.Command, args []string) error {
		stdinBytes, err := io.ReadAll(os.Stdin)
		if err != nil {
			stdinBytes = []byte{}
		}

		if tee != "" {
			if err := os.WriteFile(tee, stdinBytes, 0644); err != nil {
				fmt.Fprintf(os.Stderr, "ccglow: failed to write tee file: %v\n", err)
			}
		}

		output := run(presetName, configPath, format, string(stdinBytes))
		if output != "" {
			fmt.Print(output)
		}

		return nil
	},
}

func init() {
	rootCmd.Flags().StringVar(&presetName, "preset", "default", "Use a named preset (default, minimal, full)")
	rootCmd.Flags().StringVar(&configPath, "config", "", "Load JSON config file")
	rootCmd.Flags().StringVar(&format, "format", "ansi", "Output format: ansi, plain")
	rootCmd.Flags().StringVar(&tee, "tee", "", "Write raw stdin JSON to file before processing")

	rootCmd.SetVersionTemplate("{{.Version}}\n")
	rootCmd.SilenceUsage = true
}

// SetVersion sets the version string on the root command.
func SetVersion(v string) {
	rootCmd.Version = v
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// Version returns the build version from debug info or "dev".
func Version() string {
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}
	return "dev"
}

func run(presetName, configPath, format, stdin string) string {
	sess := session.Parse(stdin)
	if sess == nil {
		return ""
	}

	if format == "plain" {
		style.SetColorLevel(0)
	} else {
		style.SetColorLevel(1)
	}
	defer style.SetColorLevel(1)

	segments := segment.NewRegistry()
	segment.RegisterBuiltin(segments)

	providers := provider.NewRegistry()
	provider.RegisterBuiltin(providers)

	tagIdx, err := render.BuildTagIndex(providers.All())
	if err != nil {
		fmt.Fprintf(os.Stderr, "ccglow: tag index error: %v\n", err)
		return ""
	}

	tree := resolveTree(presetName, configPath)

	providerNames := render.CollectProviderNames(tree, tagIdx)
	providerData := render.ResolveProviders(providerNames, providers.All(), sess)
	segmentValues := render.ResolveSegmentValues(tagIdx, providerData)

	return render.Tree(tree, segments, sess, providerData, segmentValues, tagIdx)
}

func resolveTree(presetName, configPath string) []types.SegmentNode {
	if configPath != "" {
		data, err := os.ReadFile(configPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ccglow: failed to load config: %v\n", err)
		} else {
			tree, err := config.Parse(data)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ccglow: failed to parse config: %v\n", err)
			} else if len(tree) > 0 {
				return tree
			}
		}
	}

	if tree := preset.Get(presetName); tree != nil {
		return tree
	}

	return preset.Get("default")
}
