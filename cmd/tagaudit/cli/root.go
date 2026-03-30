package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/emm5317/tagaudit"
	"github.com/emm5317/tagaudit/internal/configfile"
	"github.com/emm5317/tagaudit/rules"
	"github.com/spf13/cobra"
)

var (
	flagConfig       string
	flagOutput       string
	flagFix          bool
	flagRules        string
	flagDisableRules string
	flagMinSeverity  string
)

var rootCmd = &cobra.Command{
	Use:   "tagaudit [flags] [packages]",
	Short: "Struct tag validation for Go",
	Long: `tagaudit validates struct tags in Go source code.

It checks for syntax errors, naming convention violations, unknown tag keys,
duplicate values, shadowed embedded tags, and more.

Packages use the same patterns as "go build" (e.g., ./..., ./internal/models).
If no packages are specified, ./... is used.`,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE:          runAnalysis,
}

func init() {
	rootCmd.Flags().StringVarP(&flagConfig, "config", "c", "", "path to config file (YAML)")
	rootCmd.Flags().StringVarP(&flagOutput, "output", "o", "text", "output format: text or json")
	rootCmd.Flags().BoolVar(&flagFix, "fix", false, "apply suggested fixes to source files")
	rootCmd.Flags().StringVar(&flagRules, "rules", "", "comma-separated rule IDs to enable")
	rootCmd.Flags().StringVar(&flagDisableRules, "disable-rules", "", "comma-separated rule IDs to disable")
	rootCmd.Flags().StringVar(&flagMinSeverity, "min-severity", "", "minimum severity: error, warning, or info")
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

func runAnalysis(cmd *cobra.Command, args []string) error {
	cfg, err := buildConfig()
	if err != nil {
		return err
	}

	patterns := args
	if len(patterns) == 0 {
		patterns = []string{"./..."}
	}

	a := tagaudit.New(cfg)
	findings, err := a.AnalyzePackages(patterns...)
	if err != nil {
		return fmt.Errorf("analysis failed: %w", err)
	}

	if flagFix {
		applied, fixErr := applyFixes(findings)
		if fixErr != nil {
			fmt.Fprintf(os.Stderr, "warning: fix failed: %v\n", fixErr)
		} else if applied > 0 {
			fmt.Fprintf(os.Stderr, "Applied %d fix(es)\n", applied)
			// Re-run analysis to show remaining findings
			findings, err = a.AnalyzePackages(patterns...)
			if err != nil {
				return fmt.Errorf("re-analysis after fix failed: %w", err)
			}
		}
	}

	switch flagOutput {
	case "json":
		return outputJSON(os.Stdout, findings)
	case "text":
		return outputText(os.Stdout, findings)
	default:
		return fmt.Errorf("unknown output format %q: must be text or json", flagOutput)
	}
}

func buildConfig() (*tagaudit.Config, error) {
	allRules := rules.All()
	var cfg *tagaudit.Config

	if flagConfig != "" {
		fc, err := configfile.Load(flagConfig)
		if err != nil {
			return nil, err
		}
		cfg, err = fc.ToConfig(allRules)
		if err != nil {
			return nil, err
		}
	} else {
		cfg = rules.DefaultConfig()
	}

	// Apply flag overrides
	if flagRules != "" || flagDisableRules != "" {
		var enable, disable []string
		if flagRules != "" {
			enable = strings.Split(flagRules, ",")
		}
		if flagDisableRules != "" {
			disable = strings.Split(flagDisableRules, ",")
		}
		cfg.Rules = configfile.FilterRulesByID(allRules, enable, disable)
	}

	if flagMinSeverity != "" {
		sev, err := configfile.ParseSeverity(flagMinSeverity)
		if err != nil {
			return nil, err
		}
		cfg.MinSeverity = tagaudit.SeverityPtr(sev)
	}

	return cfg, nil
}
