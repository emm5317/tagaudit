package severityfilter

// BadSyntax has an error-level finding (syntax) that should survive error-only filtering.
type BadSyntax struct {
	Name string `json:"name"`
	Bad  string `json:bad_no_quotes` // want `\[syntax\]`
}
