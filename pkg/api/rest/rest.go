package rest

import (
	"context"
	"fmt"
	"net/http"
	"sync"

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

func New(ctx context.Context, conn *storage.Database, cfg *config.KafkaConfig, wg *sync.WaitGroup) *Server {
	svr := &Server{
		DBConnector: conn,
		kafkaCfg:    cfg,
	}
	svr.startProducer(ctx, wg)
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

func (s *Server) startProducer(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()

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
	Data   any     `json:"data"`
	Errors []error `json:"errors"`
}

// handleGetResponse for a CRUD operation handle the response.
func handleGetResponse(w http.ResponseWriter, r *http.Request, object any, err error) {
	resp := Response{
		Data:   object,
		Errors: []error{err},
	}
	if err != nil {
		fmt.Println(err.Error())
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, resp)
		return
	}

	if object == nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, resp)
		return
	}
	render.JSON(w, r, resp)
}
