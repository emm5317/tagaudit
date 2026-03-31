package anon

// Anonymous struct in composite literal — exercises CompositeLit branch in runAnalysis.
var compLit = struct {
	UserName string `json:"userName"` // want `\[naming\]`
	Email    string `json:"email"`
}{}

// Anonymous struct in var declaration — exercises ValueSpec branch in runAnalysis.
var varDecl struct {
	Name string `json:"name"`
	Bad  string `json:bad_no_quotes` // want `\[syntax\]`
}
