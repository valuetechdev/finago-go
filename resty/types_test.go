package resty

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMultilineString(t *testing.T) {
	tests := []struct {
		Name   string
		Input  string
		Expect string
	}{
		{
			Name:   "single string",
			Input:  `"address 1"`,
			Expect: "address 1",
		},
		{
			Name:   "single multiline",
			Input:  `["address 1", "gate 44"]`,
			Expect: "address 1\ngate 44",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			var ms MultilineString

			err := json.Unmarshal([]byte(tt.Input), &ms)
			require.NoError(t, err)

			require.Equal(t, tt.Expect, ms.String())
		})
	}
}
