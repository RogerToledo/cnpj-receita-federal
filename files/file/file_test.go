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
		{
			description: "Retuns the hash of a line",
			input:       "24770086;0001;08;1;RCPM;02;20160510;00;;;20160510;2621300;4541203,4541206,4619200,4753900,6204000,7020400,7112000,7739099,8020001,8211300;RUA;PADRE CAMARGO;341;;CENTRO;84130000;PR;7735;00;11111111;;;;;;;",
			expected:    "e7e59d9c340aa2b66d44167ea1ad38a7c4b4a4e8b27f6fd3fbcc5444339ff8e0",
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