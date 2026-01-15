package httppack

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

type HTTPServer struct {
	httpHandlers *HTTPHandlers
}

func NewHTTPServer(httpHandlers HTTPHandlers) *HTTPServer {
	return &HTTPServer{
		httpHandlers: &httpHandlers,
	}
}

func (h *HTTPServer) ConnectServer() error {
	router := mux.NewRouter()

	router.Path("/employees").Methods("POST").HandlerFunc(h.httpHandlers.CreateWorkerHandler)
	router.Path("/employees").Methods("GET").HandlerFunc(h.httpHandlers.GetWorkerHandler)
	router.Path("/employees/{id}").Methods("DELETE").HandlerFunc(h.httpHandlers.DeleteWorkerHandler)

	if err := http.ListenAndServe(":5050", router); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
	return nil
}
