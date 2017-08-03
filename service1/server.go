package service1

import (
	"log"
	"net/http"

	opentracing "github.com/opentracing/opentracing-go"
)

// Do you think it's secure enough?
const Password = "opensesame"

type server struct {
	password string
}

func (s *server) serviceHandler(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	password := values.Get("passwd")
	span := opentracing.SpanFromContext(r.Context())

	if password != s.password {
		span.SetTag("error", "Invalid Password")
		http.Error(w, "Invalid Password", http.StatusBadRequest)
		return
	}

	span.LogEvent("Password verification successful")

	w.WriteHeader(http.StatusNoContent)
}

func NewServer(password string) *server {
	return &server{
		password: password,
	}
}

func (s *server) Start() {
	mux := http.NewServeMux()
	mux.Handle("/password", http.HandlerFunc(s.serviceHandler))
	if err := http.ListenAndServe(":4048", mux); err != nil {
		log.Fatalf("Error: %v\n", err)
	}
}
