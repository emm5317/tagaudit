package anon

// Anonymous struct in a composite literal — naming rule should flag userName.
var compLit = struct {
	UserName string `json:"userName"`
	Email    string `json:"email"`
}{}

// Anonymous struct in a var declaration — syntax rule should flag the bad tag.
var varDecl struct {
	Name string `json:"name"`
	Bad  string `json:bad_no_quotes`
}
