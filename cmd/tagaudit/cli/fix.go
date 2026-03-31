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

	// Check if all fixes have byte offsets for precise replacement
	allHaveOffsets := true
	for _, f := range findings {
		if f.Fix == nil {
			continue
		}
		if f.Fix.TagStart <= 0 || f.Fix.TagEnd <= 0 || f.Fix.TagEnd <= f.Fix.TagStart {
			allHaveOffsets = false
			break
		}
	}

	var applied int
	var result []byte

	if allHaveOffsets {
		result, applied = applyByteOffsetFixes(src, findings)
	} else {
		result, applied = applyLineFixes(src, findings)
	}

	if applied == 0 {
		return 0, nil
	}

	// Format with go/format to ensure valid Go
	formatted, err := format.Source(result)
	if err != nil {
		// Write unformatted if formatting fails
		formatted = result
	}

	if err := os.WriteFile(path, formatted, 0644); err != nil {
		return 0, err
	}

	return applied, nil
}

// applyByteOffsetFixes uses exact byte offsets to replace tag literals.
func applyByteOffsetFixes(src []byte, findings []tagaudit.Finding) ([]byte, int) {
	// Sort by TagStart descending so replacements don't shift earlier offsets
	sort.Slice(findings, func(i, j int) bool {
		return findings[i].Fix.TagStart > findings[j].Fix.TagStart
	})

	applied := 0
	result := make([]byte, len(src))
	copy(result, src)

	for _, f := range findings {
		if f.Fix == nil {
			continue
		}
		start := f.Fix.TagStart
		end := f.Fix.TagEnd
		if start <= 0 || end <= 0 || end > len(result) || start >= end {
			continue
		}

		replacement := []byte("`" + f.Fix.NewTagValue + "`")
		result = append(result[:start], append(replacement, result[end:]...)...)
		applied++
	}

	return result, applied
}

// applyLineFixes is the fallback that uses line-based backtick search.
func applyLineFixes(src []byte, findings []tagaudit.Finding) ([]byte, int) {
	lines := strings.Split(string(src), "\n")

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

	return []byte(strings.Join(lines, "\n")), applied
}
