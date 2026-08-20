package payday

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// requireOnline skips tests that talk to the live Payday API. `go test -short`
// must pass offline and without credentials.
func requireOnline(t *testing.T) {
	t.Helper()
	if testing.Short() {
		t.Skip("skipping test that requires the Payday API")
	}
}

func TestClient(t *testing.T) {
	requireOnline(t)
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
