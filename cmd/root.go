package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime/debug"

	"github.com/jheddings/ccglow/internal/config"
	"github.com/jheddings/ccglow/internal/preset"
	"github.com/jheddings/ccglow/internal/provider"
	"github.com/jheddings/ccglow/internal/render"
	"github.com/jheddings/ccglow/internal/session"
	"github.com/jheddings/ccglow/internal/style"
	"github.com/jheddings/ccglow/internal/types"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	presetName string
	configPath string
	format     string
	tee     string
	dump    string
	logPath    string
	verbose    bool
)

var rootCmd = &cobra.Command{
	Use:   "ccglow",
	Short: "Composable statusline for Claude Code",
	Long:  "Reads session JSON from stdin, outputs styled statusline to stdout.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if logPath == "" {
			log.Logger = zerolog.Nop()
		} else {
			f, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				return fmt.Errorf("failed to open log file: %w", err)
			}
			log.Logger = zerolog.New(f).With().Timestamp().Logger()
		}

		if verbose {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		} else {
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		stdinBytes, err := io.ReadAll(os.Stdin)
		if err != nil {
			stdinBytes = []byte{}
		}

		if tee != "" {
			if err := os.WriteFile(tee, stdinBytes, 0644); err != nil {
				log.Error().Err(err).Msg("failed to write tee file")
			}
		}

		output, env := run(presetName, configPath, format, string(stdinBytes))

		if dump != "" {
			data, err := json.MarshalIndent(env, "", "  ")
			if err != nil {
				log.Error().Err(err).Msg("failed to marshal env")
			} else if err := os.WriteFile(dump, data, 0644); err != nil {
				log.Error().Err(err).Msg("failed to write dump file")
			}
		}

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
	rootCmd.Flags().StringVar(&dump, "dump", "", "Write resolved provider env as JSON to file")

	rootCmd.PersistentFlags().StringVar(&logPath, "log", "", "Write logs to file (no logging when omitted)")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Set log level to debug")

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

func run(presetName, configPath, format, stdin string) (string, map[string]any) {
	sess := session.Parse(stdin)
	if sess == nil {
		log.Warn().Msg("failed to parse session")
		return "", nil
	}
	log.Debug().Str("cwd", sess.CWD).Msg("session parsed")

	if format == "plain" {
		style.SetColorLevel(0)
	} else {
		style.SetColorLevel(1)
	}
	defer style.SetColorLevel(1)

	providers := provider.NewRegistry()
	provider.RegisterBuiltin(providers)

	env, defaultFormats := render.BuildEnv(providers.All(), sess)
	log.Debug().Int("providers", len(env)).Msg("env built")

	tree := resolveTree(presetName, configPath)
	log.Debug().Int("count", len(tree)).Msg("tree resolved")

	output := render.Tree(tree, sess, env, defaultFormats)
	log.Debug().Msg("render complete")

	return output, env
}

func resolveTree(presetName, configPath string) []types.SegmentNode {
	if configPath != "" {
		data, err := os.ReadFile(configPath)
		if err != nil {
			log.Error().Err(err).Str("path", configPath).Msg("failed to load config")
		} else {
			log.Debug().Int("bytes", len(data)).Str("path", configPath).Msg("config file read")
			tree, err := config.Parse(data)
			if err != nil {
				log.Error().Err(err).Str("path", configPath).Msg("failed to parse config")
			} else if len(tree) > 0 {
				log.Debug().Int("count", len(tree)).Str("path", configPath).Msg("config tree parsed")
				return tree
			} else {
				log.Warn().Str("path", configPath).Msg("config file produced empty tree")
			}
		}
	}

	if tree := preset.Get(presetName); tree != nil {
		log.Debug().Str("preset", presetName).Msg("using preset")
		return tree
	}

	log.Debug().Msg("falling back to default preset")
	return preset.Get("default")
}
