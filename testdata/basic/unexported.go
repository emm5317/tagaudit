package basic

type WithUnexported struct {
	Name    string `json:"name"`    // ok - exported
	email   string `json:"email"`   // bad - unexported with json tag
	secret  string `json:"secret"`  // bad - unexported with json tag
	private int    // ok - no tag
}
