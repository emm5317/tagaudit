package naming

type User struct {
	UserName  string `json:"userName"`  // should be user_name
	Email     string `json:"email"`     // ok
	CreatedAt string `json:"createdAt"` // should be created_at
	Age       int    `json:"age"`       // ok
}

type CamelCaseOK struct {
	FirstName string `json:"first_name"` // ok for snake_case
}
