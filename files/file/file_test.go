package file

import "testing"

func TestGetLineHash(t *testing.T) {
	cases := []struct {
		description string
		input       string
		expected    string
	}{
		{
			description: "Retuns the hash of a line empty",
			input:       "",
			expected:    "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
	}

	for _, tc := range cases {
		t.Log(tc.description)
		output := getLineHash(tc.input)
		if output != tc.expected {
			t.Errorf("getLineHash(%q) = %q; want %q", tc.input, output, tc.expected)
		}
	}
}
