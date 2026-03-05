package intake

import "net/http"

const EventsPath = "/v1/silver/events"

func NewMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle(EventsPath, NewHandler())
	return mux
}
