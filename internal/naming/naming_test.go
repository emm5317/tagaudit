package naming

import "testing"

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"UserName", "user_name"},
		{"userName", "user_name"},
		{"user_name", "user_name"},
		{"HTTPServer", "http_server"},
		{"userID", "user_id"},
		{"APIKey", "api_key"},
		{"HTMLParser", "html_parser"},
		{"simpleTest", "simple_test"},
		{"JSONResponse", "json_response"},
		{"getHTTPResponse", "get_http_response"},
		{"ID", "id"},
		{"", ""},
		{"a", "a"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := ToSnakeCase(tt.input)
			if got != tt.want {
				t.Errorf("ToSnakeCase(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestToCamelCase(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"user_name", "userName"},
		{"UserName", "userName"},
		{"http_server", "httpServer"},
		{"api_key", "apiKey"},
		{"simple", "simple"},
		{"", ""},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := ToCamelCase(tt.input)
			if got != tt.want {
				t.Errorf("ToCamelCase(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestToPascalCase(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"user_name", "UserName"},
		{"userName", "UserName"},
		{"http_server", "HTTPServer"},
		{"", ""},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := ToPascalCase(tt.input)
			if got != tt.want {
				t.Errorf("ToPascalCase(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestToKebabCase(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"UserName", "user-name"},
		{"userName", "user-name"},
		{"user-name", "user-name"},
		{"HTTPServer", "http-server"},
		{"", ""},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := ToKebabCase(tt.input)
			if got != tt.want {
				t.Errorf("ToKebabCase(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestMatchesConvention(t *testing.T) {
	tests := []struct {
		input      string
		convention string
		want       bool
	}{
		{"user_name", "snake_case", true},
		{"userName", "snake_case", false},
		{"userName", "camelCase", true},
		{"UserName", "camelCase", false},
		{"UserName", "PascalCase", true},
		{"user-name", "kebab-case", true},
		{"user_name", "kebab-case", false},
	}
	for _, tt := range tests {
		t.Run(tt.input+"_"+tt.convention, func(t *testing.T) {
			got := MatchesConvention(tt.input, tt.convention)
			if got != tt.want {
				t.Errorf("MatchesConvention(%q, %q) = %v, want %v", tt.input, tt.convention, got, tt.want)
			}
		})
	}
}

func TestConvert_UnknownConvention(t *testing.T) {
	// Unknown convention should return input unchanged
	got := Convert("fooBar", "SCREAMING_CASE")
	if got != "fooBar" {
		t.Errorf("Convert with unknown convention should return input, got %q", got)
	}
}

func TestTitleWord_Empty(t *testing.T) {
	got := titleWord("")
	if got != "" {
		t.Errorf("titleWord(\"\") = %q, want \"\"", got)
	}
}
