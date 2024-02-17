// Portions of this file were taken from https://github.com/tonyghita/graphql-go-example/blob/main/schema/schema.go
package schema

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"strings"

	"github.com/graph-gophers/graphql-go"
	"github.com/sassoftware/event-provenance-registry/pkg/api/graphql/resolvers"
	"github.com/sassoftware/event-provenance-registry/pkg/message"
	"github.com/sassoftware/event-provenance-registry/pkg/storage"
)

func New(connection *storage.Database, msgProducer message.TopicProducer) *graphql.Schema {
	s, err := String()
	if err != nil {
		log.Fatalf("reading embedded schema contents: %s", err)
	}
	opts := []graphql.SchemaOpt{graphql.UseFieldResolvers()}
	return graphql.MustParseSchema(s, resolvers.New(connection, msgProducer), opts...)
}

//go:embed *.graphql types/*.graphql
var content embed.FS

// TAKEN FROM GRAPHQL EXAMPLE CODE AT https://github.com/tonyghita/graphql-go-example/blob/main/schema/schema.go
func String() (string, error) {
	var buf bytes.Buffer

	fn := func(path string, _ fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walking dir: %w", err)
		}

		// Only add files with the .graphql extension.
		if !strings.HasSuffix(path, ".graphql") {
			return nil
		}

		b, err := content.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading file %q: %w", path, err)
		}

		// Add a newline to separate each file.
		b = append(b, []byte("\n")...)

		if _, err := buf.Write(b); err != nil {
			return fmt.Errorf("writing %q bytes to buffer: %w", path, err)
		}

		return nil
	}

	// Recursively walk this directory and append all the file contents together.
	if err := fs.WalkDir(content, ".", fn); err != nil {
		return buf.String(), fmt.Errorf("walking content directory: %w", err)
	}

	return buf.String(), nil
}
