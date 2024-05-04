package rest

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	eprErrors "github.com/sassoftware/event-provenance-registry/pkg/errors"
	"github.com/sassoftware/event-provenance-registry/pkg/message"
	"github.com/sassoftware/event-provenance-registry/pkg/storage"
)

type Server struct {
	DBConnector *storage.Database

	msgProducer message.TopicProducer
}

func New(conn *storage.Database, msgProducer message.TopicProducer) *Server {
	svr := &Server{
		DBConnector: conn,
		msgProducer: msgProducer,
	}
	return svr
}

func (s *Server) ServeOpenAPIDoc(_ string) http.HandlerFunc {
	return func(_ http.ResponseWriter, _ *http.Request) {
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

func handleResponse(w http.ResponseWriter, r *http.Request, data any, err error) {
	resp := Response{
		Data: data,
	}
	if err == nil {
		render.Status(r, http.StatusOK)
		render.JSON(w, r, resp)
		return
	}

	var status int
	switch err.(type) {
	case eprErrors.MissingObjectError:
		status = http.StatusNotFound
	case eprErrors.InvalidInputError:
		status = http.StatusBadRequest
	default:
		status = http.StatusInternalServerError
	}
	render.Status(r, status)

	if status == http.StatusInternalServerError {
		// don't expose server internals
		resp.Errors = []string{"internal server error"}
	} else {
		resp.Errors = []string{err.Error()}
	}

	slog.Error("error during request", "error", err, "url", r.URL)
	render.JSON(w, r, resp)
}
