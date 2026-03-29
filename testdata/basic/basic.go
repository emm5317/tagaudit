package basic

// GoodStruct has well-formed tags.
type GoodStruct struct {
	Name  string `json:"name"`
	Email string `json:"email" db:"email"`
	Age   int    `json:"age,omitempty"`
}

// BadSyntax has a malformed tag.
type BadSyntax struct {
	Name string `json:"name"`
	Bad  string `json:bad_no_quotes`
}

// NoTags has no tags at all.
type NoTags struct {
	Name string
	Age  int
}
