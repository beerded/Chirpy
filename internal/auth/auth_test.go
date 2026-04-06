package auth

import "testing"

func TestSmoketest(t *testing.T) {
	password1 := "correctPassword"
	password2 := "anotherCorrectPassword"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	cases := []struct {
		testName	string
		password	string
		hash		string
		wantErr		bool
		wantMatch	bool
	}{
		{
			testName:		"Correct Password",
			password:		password1,
			hash:			hash1,
			wantErr:		false,
			wantMatch:		true,
		},
		{
			testName:		"Incorrect Password",
			password:		"Not the right password",
			hash:			hash1,
			wantErr:		false,
			wantMatch:		false,
		},
		{
			testName:		"Other Password",
			password:		password1,
			hash:			hash2,
			wantErr:		false,
			wantMatch:		false,
		},
		{
			testName:		"Empty Password",
			password:		"",
			hash:			hash1,
			wantErr:		false,
			wantMatch:		false,
		},
		{
			testName:		"Invalid Hash",
			password:		password1,
			hash:			"invalid hash",
			wantErr:		true,
			wantMatch:		false,
		},
	}

	for _, c := range cases {
		ok, err := CheckPasswordHash(c.password, c.hash)
		if (err != nil) != c.wantErr {
			t.Errorf("TEST %v - CheckPasswordHash() had err = %v, want err = %v", c.testName, err, c.wantErr)
		}
		if !c.wantErr && ok != c.wantMatch {
			t.Errorf("TEST %v - CheckPasswordHash() had match = %v, want match = %v", c.testName, ok, c.wantMatch)
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
