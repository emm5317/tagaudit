package basic

type Good struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type BadSyntax struct {
	Name string `json:"name"`
	Bad  string `json:bad_no_quotes` // want `\[syntax\]` `\[completeness\]`
}

type BadNaming struct {
	UserName string `json:"userName"` // want `\[naming\]`
}

type WithUnexported struct {
	Name  string `json:"name"`
	email string `json:"email"` // want `\[unexported\]`
}

// Non-struct type declarations — exercise the early-return branches
type MyString string

type MyInterface interface {
	Foo()
}

type MyAlias = string

// Incomplete struct — triggers completeness (struct-level check)
type Incomplete struct {
	Name string `json:"name"`
	Age  int    // want `\[completeness\]`
}

// Embedded struct for shadow/duplicate coverage
type Base struct {
	ID int `json:"id"`
}

type WithShadow struct {
	Base
	Title string `json:"id"` // want `\[shadow\]` `\[duplicates\]`
}
