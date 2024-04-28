package middleware

import (
	"github.com/Atvit/assessment-tax/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAuthenticate(t *testing.T) {
	type testcase struct {
		Name        string
		Username    string
		Password    string
		Expected    bool
		ExpectedErr error
	}

	cfg := config.Configuration{AdminUsername: "username", AdminPassword: "p@ssw0rd"}
	tcs := []testcase{
		{"username matches but password", "username", "password", false, nil},
		{"password matches but username", "uname", "p@ssw0rd", false, nil},
		{"username and password do not match", "uname", "pwd", false, nil},
		{"matches both the username and password", "username", "p@ssw0rd", true, nil},
	}

	for _, tc := range tcs {
		t.Run(tc.Name, func(t *testing.T) {
			match, err := Authenticate(tc.Username, tc.Password, &cfg)

			assert.Equal(t, tc.ExpectedErr, err)
			assert.Equal(t, tc.Expected, match)
		})
	}
}
