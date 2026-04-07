package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestLogin(t *testing.T) {
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
func TestJWTS(t *testing.T) {
	secretString1 := "foobar"
	secretString2 := "barfoo"
	uuid1, _ := uuid.NewRandom()
	uuid2, _ := uuid.NewRandom()
	token1, _ := MakeJWT(uuid1, secretString1, 10 * time.Minute)
	token2, _ := MakeJWT(uuid2, secretString2, -10 * time.Minute)
	t.Logf("token1: %s", token1)
	t.Logf("token2: %s", token2)

	cases := []struct {
		testName	string
		secret		string
		uuid		uuid.UUID
		token		string
		wantUUID	uuid.UUID
		wantErr		bool
	}{
		{
			testName:		"Correct Secret",
			secret:			secretString1,
			uuid:			uuid1,
			token:			token1,
			wantUUID:		uuid1,
			wantErr:		false,
		},
		{
			testName:		"Incorrect Secret",
			secret:			"wrong secret",
			uuid:			uuid1,
			token:			token1,
			wantUUID:		uuid.Nil,
			wantErr:		true,
		},
		{
			testName:		"Expired Token",
			secret:			secretString2,
			uuid:			uuid2,
			token:			token2,
			wantUUID:		uuid.Nil,
			wantErr:		true,
		},
		{
			testName:		"Invalid Token",
			secret:			secretString1,
			uuid:			uuid1,
			token:			"invalid.token",
			wantUUID:		uuid.Nil,
			wantErr:		true,
		},
	}

	for _, c := range cases {
		ss, err := ValidateJWT(c.token, c.secret)
		if (err != nil) != c.wantErr {
			t.Errorf("TEST %v - ValidateJWT() had err = %v, want err = %v\n", c.testName, err, c.wantErr)
			return
		}
		if ss != c.wantUUID {
			t.Errorf("TEST %v - ValidateJWT() wanted UUID %v, got %v\n", c.testName, c.wantUUID, ss)
		}
	}
}
