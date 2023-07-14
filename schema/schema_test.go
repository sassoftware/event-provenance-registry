package schema_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.sas.com/async-event-infrastructure/server/schema"
)

func TestString(t *testing.T) {
	s, err := schema.String()

	require.NoError(t, err)
	require.NotEmpty(t, s)
}
