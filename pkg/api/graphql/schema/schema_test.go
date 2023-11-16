package schema_test

import (
	"testing"

	"github.com/sassoftware/event-provenance-registry/pkg/api/graphql/schema"
	"github.com/stretchr/testify/require"
)

func TestString(t *testing.T) {
	s, err := schema.String()

	require.NoError(t, err)
	require.NotEmpty(t, s)
}
