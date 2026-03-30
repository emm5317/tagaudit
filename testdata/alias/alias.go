package alias

// Base is a plain struct with tags.
type Base struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// AliasedStruct is a type alias for Base. tagaudit should resolve
// through the alias and check the underlying struct's fields.
type AliasedStruct = Base

// UsesAlias embeds the alias and adds a field with a naming violation.
type UsesAlias struct {
	AliasedStruct
	UserName string `json:"userName"` // camelCase — naming rule should flag
}
