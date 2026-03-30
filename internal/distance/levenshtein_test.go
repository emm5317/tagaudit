package distance

import "testing"

func TestLevenshtein(t *testing.T) {
	tests := []struct {
		a, b string
		want int
	}{
		{"", "", 0},
		{"abc", "", 3},
		{"", "abc", 3},
		{"abc", "abc", 0},
		{"json", "josn", 2},
		{"omitempty", "omitemtpy", 2},
		{"json", "yaml", 4},
		{"a", "b", 1},
		{"kitten", "sitting", 3},
	}
	for _, tt := range tests {
		got := Levenshtein(tt.a, tt.b)
		if got != tt.want {
			t.Errorf("Levenshtein(%q, %q) = %d, want %d", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestClosestMatch(t *testing.T) {
	candidates := []string{"json", "yaml", "xml", "toml", "db"}

	tests := []struct {
		input   string
		maxDist int
		want    string
		wantOK  bool
	}{
		{"josn", 2, "json", true},
		{"ymal", 2, "yaml", true},
		{"xm", 2, "xml", true},
		{"totally_different", 2, "", false},
		{"json", 0, "json", true},
	}
	for _, tt := range tests {
		got, ok := ClosestMatch(tt.input, candidates, tt.maxDist)
		if ok != tt.wantOK || got != tt.want {
			t.Errorf("ClosestMatch(%q, candidates, %d) = (%q, %v), want (%q, %v)",
				tt.input, tt.maxDist, got, ok, tt.want, tt.wantOK)
		}
	}
}
