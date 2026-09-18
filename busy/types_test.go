package busy

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEmailUnmarshal(t *testing.T) {
	tests := []struct {
		Name      string
		Input     string
		Expect    string
		ExpectErr bool
	}{
		{
			Name:   "plain address",
			Input:  `"user@example.org"`,
			Expect: "user@example.org",
		},
		{
			// Busy returns this for users without an e-mail address.
			Name:   "empty address",
			Input:  `""`,
			Expect: "",
		},
		{
			Name:   "address with display name is normalized",
			Input:  `"John Doe <user@example.org>"`,
			Expect: "user@example.org",
		},
		{
			Name:      "not an address",
			Input:     `"not-an-email"`,
			ExpectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			var e Email

			err := json.Unmarshal([]byte(tt.Input), &e)
			if tt.ExpectErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.Expect, e.String())
			require.Equal(t, tt.Expect == "", e.IsEmpty())
		})
	}
}

func TestEmailMarshal(t *testing.T) {
	tests := []struct {
		Name      string
		Input     Email
		Expect    string
		ExpectErr bool
	}{
		{
			Name:   "plain address",
			Input:  Email("user@example.org"),
			Expect: `"user@example.org"`,
		},
		{
			Name:   "empty address",
			Input:  Email(""),
			Expect: `""`,
		},
		{
			Name:      "not an address",
			Input:     Email("not-an-email"),
			ExpectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			b, err := json.Marshal(tt.Input)
			if tt.ExpectErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.JSONEq(t, tt.Expect, string(b))
		})
	}
}

// TestEmailInGeneratedStruct pins the behaviour the overlay exists for: a user
// payload with an empty e-mail must decode without error.
func TestEmailInGeneratedStruct(t *testing.T) {
	require := require.New(t)

	var u User
	require.NoError(json.Unmarshal([]byte(`{"email":""}`), &u))
	require.True(u.Email.IsEmpty())

	require.NoError(json.Unmarshal([]byte(`{"email":"user@example.org"}`), &u))
	require.Equal("user@example.org", u.Email.String())
}
