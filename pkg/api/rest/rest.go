package rest

import (
	"context"
	"fmt"
	"gitlab.sas.com/async-event-infrastructure/server/pkg/config"
	"net/http"
	"sync"

	"github.com/go-chi/render"
	"gitlab.sas.com/async-event-infrastructure/server/pkg/storage"
)

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

// RestResponse generic rest response for all object types.
type RestResponse struct {
	Data   any     `json:"data"`
	Errors []error `json:"errors"`
}

// handleGetResponse for a CRUD operation handle the response.
func handleGetResponse(w http.ResponseWriter, r *http.Request, object any, err error) {
	resp := RestResponse{
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
