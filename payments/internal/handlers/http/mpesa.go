package http

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jwambugu/mpesa-golang-sdk"
)

func (s *HTTPServer) registerCallbackRoutes(r *chi.Mux) {
	r.Get("/callback", s.handleMpesaCallback)
}

func (s *HTTPServer) handleMpesaCallback(w http.ResponseWriter, r *http.Request) {
	callback, err := mpesa.UnmarshalSTKPushCallback(r)
	if err != nil {
		log.Fatalln(err)
	}

	err = s.PaymentsService.HandleMpesaCallback(r.Context(), &callback.Body.STKCallback)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("%+v", callback)
}
