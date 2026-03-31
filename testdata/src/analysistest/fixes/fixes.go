package fixes

// NamingFix exercises the buildDiagnostic suggested-fix path.
type NamingFix struct {
	UserName string `json:"userName"` // want `\[naming\]`
}

// Incomplete exercises struct-level check field matching in runAnalysis.
type Incomplete struct {
	Name string `json:"name"`
	Age  int    // want `\[completeness\]`
}

// WithShadow exercises struct-level field name matching for shadow/duplicates.
type Base struct {
	ID int `json:"id"`
}

type WithShadow struct {
	Base
	Title string `json:"id"` // want `\[shadow\]` `\[duplicates\]`
}
