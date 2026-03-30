package cli

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/emm5317/tagaudit"
	"go/format"
)

// applyFixes writes suggested fixes to source files. Returns the number
// of fixes applied.
func applyFixes(findings []tagaudit.Finding) (int, error) {
	// Group fixable findings by file
	byFile := make(map[string][]tagaudit.Finding)
	for _, f := range findings {
		if f.Fix == nil || f.Pos.Filename == "" {
			continue
		}
		byFile[f.Pos.Filename] = append(byFile[f.Pos.Filename], f)
	}

	applied := 0
	for file, fixes := range byFile {
		n, err := applyFixesToFile(file, fixes)
		if err != nil {
			return applied, fmt.Errorf("fixing %s: %w", file, err)
		}
		applied += n
	}

	return applied, nil
}

func applyFixesToFile(path string, findings []tagaudit.Finding) (int, error) {
	src, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}

	content := string(src)
	lines := strings.Split(content, "\n")

	// Sort findings by line in reverse order so replacements don't shift offsets
	sort.Slice(findings, func(i, j int) bool {
		return findings[i].Pos.Line > findings[j].Pos.Line
	})

	applied := 0
	for _, f := range findings {
		if f.Fix == nil || f.Pos.Line <= 0 || f.Pos.Line > len(lines) {
			continue
		}

		line := lines[f.Pos.Line-1]
		// Find the backtick-quoted tag in this line
		tagStart := strings.Index(line, "`")
		tagEnd := strings.LastIndex(line, "`")
		if tagStart < 0 || tagEnd <= tagStart {
			continue
		}

		// Replace the tag content
		newLine := line[:tagStart+1] + f.Fix.NewTagValue + line[tagEnd:]
		lines[f.Pos.Line-1] = newLine
		applied++
	}

	if applied == 0 {
		return 0, nil
	}

	newContent := strings.Join(lines, "\n")

	// Format with go/format to ensure valid Go
	formatted, err := format.Source([]byte(newContent))
	if err != nil {
		// Write unformatted if formatting fails
		formatted = []byte(newContent)
	}

	if err := os.WriteFile(path, formatted, 0644); err != nil {
		return 0, err
	}

	return applied, nil
}
