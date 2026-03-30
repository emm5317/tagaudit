package configfile

import (
	"fmt"
	"os"
	"strings"

	"github.com/emm5317/tagaudit"
	"gopkg.in/yaml.v3"
)

// FileConfig represents the structure of a .tagaudit.yaml config file.
type FileConfig struct {
	Rules             RulesConfig         `yaml:"rules"`
	NamingConventions map[string]string   `yaml:"naming_conventions"`
	RequiredTagKeys   []string            `yaml:"required_tag_keys"`
	KnownTagKeys      []string            `yaml:"known_tag_keys"`
	KnownOptions      map[string][]string `yaml:"known_options"`
	MinSeverity       string              `yaml:"min_severity"`
}

// RulesConfig controls which rules are enabled/disabled.
type RulesConfig struct {
	Enable  []string `yaml:"enable"`
	Disable []string `yaml:"disable"`
}

// Load reads and parses a YAML config file.
func Load(path string) (*FileConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var fc FileConfig
	if err := yaml.Unmarshal(data, &fc); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	return &fc, nil
}

// ToConfig converts a FileConfig into a tagaudit.Config using the given
// set of available rules (typically rules.All()).
func (fc *FileConfig) ToConfig(allRules []tagaudit.Rule) (*tagaudit.Config, error) {
	cfg := &tagaudit.Config{
		NamingConventions: fc.NamingConventions,
		RequiredTagKeys:   fc.RequiredTagKeys,
		KnownTagKeys:      fc.KnownTagKeys,
		KnownOptions:      fc.KnownOptions,
	}

	// Filter rules based on enable/disable lists
	cfg.Rules = filterRules(allRules, fc.Rules.Enable, fc.Rules.Disable)

	// Parse severity
	if fc.MinSeverity != "" {
		sev, err := parseSeverity(fc.MinSeverity)
		if err != nil {
			return nil, err
		}
		cfg.MinSeverity = sev
	}

	return cfg, nil
}

func filterRules(allRules []tagaudit.Rule, enable, disable []string) []tagaudit.Rule {
	if len(enable) == 0 && len(disable) == 0 {
		return allRules
	}

	disableSet := make(map[string]bool, len(disable))
	for _, id := range disable {
		disableSet[id] = true
	}

	if len(enable) > 0 {
		enableSet := make(map[string]bool, len(enable))
		for _, id := range enable {
			enableSet[id] = true
		}
		var filtered []tagaudit.Rule
		for _, r := range allRules {
			if enableSet[r.ID()] && !disableSet[r.ID()] {
				filtered = append(filtered, r)
			}
		}
		return filtered
	}

	// Only disable list provided — include all except disabled
	var filtered []tagaudit.Rule
	for _, r := range allRules {
		if !disableSet[r.ID()] {
			filtered = append(filtered, r)
		}
	}
	return filtered
}

// ParseSeverity converts a severity string to a tagaudit.Severity value.
func ParseSeverity(s string) (tagaudit.Severity, error) {
	return parseSeverity(s)
}

func parseSeverity(s string) (tagaudit.Severity, error) {
	switch strings.ToLower(s) {
	case "error":
		return tagaudit.SeverityError, nil
	case "warning", "warn":
		return tagaudit.SeverityWarning, nil
	case "info":
		return tagaudit.SeverityInfo, nil
	default:
		return 0, fmt.Errorf("unknown severity %q: must be error, warning, or info", s)
	}
}

// FilterRulesByID filters rules from allRules to include only those with
// IDs in enable (if non-empty) and exclude those with IDs in disable.
func FilterRulesByID(allRules []tagaudit.Rule, enable, disable []string) []tagaudit.Rule {
	return filterRules(allRules, enable, disable)
}
