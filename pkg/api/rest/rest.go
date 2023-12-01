package rest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/sassoftware/event-provenance-registry/pkg/config"
	"github.com/sassoftware/event-provenance-registry/pkg/storage"
	"github.com/sassoftware/event-provenance-registry/pkg/utils"
)

var logger = utils.MustGetLogger("server", "pkg.api.rest")

type Server struct {
	DBConnector *storage.Database

	kafkaCfg *config.KafkaConfig
}

func New(ctx context.Context, conn *storage.Database, cfg *config.KafkaConfig) *Server {
	svr := &Server{
		DBConnector: conn,
		kafkaCfg:    cfg,
	}
	svr.startProducer(ctx)
	return svr
}

func (s *Server) ServeOpenAPIDoc(_ string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// implement using swaggest/rest instead of statically
		// writing an openapi doc
		// https://github.com/swaggest/rest/
		panic("implement me!")
	}
}

func (s *Server) startProducer(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-s.kafkaCfg.MsgChannel:
				s.kafkaCfg.Producer.Async(s.kafkaCfg.Topic, msg)
			}
		}
	}()
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
