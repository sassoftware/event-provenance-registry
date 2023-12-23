package rest

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/sassoftware/event-provenance-registry/pkg/message"
	"github.com/sassoftware/event-provenance-registry/pkg/storage"
	"github.com/sassoftware/event-provenance-registry/pkg/utils"
)

var logger = utils.MustGetLogger("server", "pkg.api.rest")

// Server type represents a server with a database connector and Kafka configuration.
// @property DBConnector - DBConnector is a pointer to a storage.Database object. It is used to
// establish a connection to a database and perform database operations.
// @property kafkaCfg - The `kafkaCfg` property is a configuration object for Kafka. It is of type
// `*config.KafkaConfig`, which means it is a pointer to an instance of the `KafkaConfig` struct
// defined in the `config` package. This object contains various configuration parameters related to
// Kafka,
type Server struct {
	DBConnector *storage.Database

	msgProducer message.TopicProducer
}

// New function creates a new Server instance and starts a Kafka producer.
func New(ctx context.Context, conn *storage.Database, cfg *config.KafkaConfig) *Server {
	svr := &Server{
		DBConnector: conn,
		msgProducer: msgProducer,
	}
	return svr
}

// ServeOpenAPIDoc function is a method of the `Server` struct. It returns an `http.HandlerFunc`
// that handles requests for serving the OpenAPI documentation.
func (s *Server) ServeOpenAPIDoc(_ string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// implement using swaggest/rest instead of statically
		// writing an openapi doc
		// https://github.com/swaggest/rest/
		panic("implement me!")
	}
}

// Response generic rest response for all object types.
type Response struct {
	Data   any      `json:"data,omitempty"`
	Errors []string `json:"errors,omitempty"`
}

type missingObjectError struct {
	msg string
}

func (m missingObjectError) Error() string {
	return fmt.Sprintf("object not found: %s", m.msg)
}

type invalidInputError struct {
	msg string
}

func (n invalidInputError) Error() string {
	return fmt.Sprintf("request parameters invalid: %s", n.msg)
}

func handleResponse(w http.ResponseWriter, r *http.Request, data any, err error) {
	resp := Response{
		Data: data,
	}
	if err == nil {
		render.Status(r, http.StatusOK)
		render.JSON(w, r, resp)
		return
	}

	resp.Errors = []string{err.Error()}
	switch err.(type) {
	case missingObjectError:
		render.Status(r, http.StatusNotFound)
	case invalidInputError:
		render.Status(r, http.StatusBadRequest)
	default:
		render.Status(r, http.StatusInternalServerError)
	}
	logger.Error(err, "error during request", "url", r.URL)
	render.JSON(w, r, resp)
}
