package embedded

// Base has some tagged fields.
type Base struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// DuplicateTag has a field whose json tag duplicates an embedded field's tag.
type DuplicateTag struct {
	Base
	Name string `json:"name"` // duplicates Base.Name's json tag
}

// ShadowedTag has a field that shadows an embedded field's tag value.
type ShadowedTag struct {
	Base
	DisplayName string `json:"name"` // shadows Base.Name's json tag "name"
}

// Complete has all exported fields tagged.
type Complete struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Incomplete has some exported fields without json tags.
type Incomplete struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    // missing json tag
	Phone string // missing json tag
}

// NoDuplicates has unique tag values.
type NoDuplicates struct {
	Base
	Email string `json:"email"` // unique
}
