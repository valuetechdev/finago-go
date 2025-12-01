package payday

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	require := require.New(t)

	c := New(os.Getenv("TFSO_PAYROLL_SECRET"))
	require.NotNil(c)

	require.False(c.IsTokenValid(), "token should be invalid before init")
	require.NoError(c.Authenticate(), "client should authenticate")
	require.True(c.IsTokenValid(), "token should be valid after authentication")

	a, err := c.GetAbsenceV2WithResponse(context.TODO(), &GetAbsenceV2Params{})
	require.NoError(err, "GetAbsenceV2EmpIdWithResponse")
	require.Empty(a.JSON200)
}
