package auth

import "testing"

func TestSmoketest(t *testing.T) {
	cases := []struct {
		input		string
		expected	bool
	}{
		{
			input:		"foobar",
			expected:	true,
		},
		{
			input:		"",
			expected:	true,
		},
	}

	for _, c := range cases {
		hash, err := HashPassword(c.input)
		if err != nil {
			t.Errorf("Got an error when hashing %s: %v", c.input, err)
			continue
		}
		actual, err := CheckPasswordHash(c.input, hash)
		if actual != c.expected {
			t.Errorf("Hashes are expected to match but they do not")
		}
	}
}

func TestNegativeCase(t *testing.T) {
	cases := []struct {
		input		string
		expected	bool
	}{
		{
			input:		"foobar",
			expected:	false,
		},
		{
			input:		"",
			expected:	false,
		},
	}

	for _, c := range cases {
		hash, err := HashPassword(c.input)
		if err != nil {
			t.Errorf("Got an error when hashing %s: %v", c.input, err)
			continue
		}
		actual, err := CheckPasswordHash("some other password", hash)
		if actual != c.expected {
			t.Errorf("Hashes are expected to be different but they somehow still match")
		}
	}
}
