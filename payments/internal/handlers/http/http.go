package http

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/leta/order-management-system/payments/internal/service"

	"golang.org/x/crypto/acme/autocert"
)

const (
	// ShutdownTimeout is the time given for outstanding requests to finish before shutdown.
	ShutdownTimeout = 1 * time.Second
	ReadTimeout     = 5 * time.Second
	WriteTimeout    = 10 * time.Second
	IdleTimeout     = 120 * time.Second
)

type HTTPServer struct {
	ln     net.Listener
	server *http.Server
	router *chi.Mux

	// Bind address & domain for the server's listener.
	// If domain is specified, server is run on TLS using acme/autocert.
	Addr   string
	Domain string

	// Services
	PaymentsService service.PaymentsService
}

// NewHTTPServer creates a new instance of HTTPServer.
func NewHTTPServer() *HTTPServer {

	s := &HTTPServer{
		router: chi.NewRouter(),
		server: &http.Server{
			ReadTimeout:  ReadTimeout,
			WriteTimeout: WriteTimeout,
			IdleTimeout:  IdleTimeout,
		},
	}

	s.router.Use(middleware.Logger)

	s.registerCallbackRoutes(s.router)

	return s
}

// UseTLS returns true if the cert & key file are specified.
func (s *HTTPServer) UseTLS() bool {
	return s.Domain != ""
}

// Scheme returns the URL scheme for the server.
func (s *HTTPServer) Scheme() string {
	if s.UseTLS() {
		return "https"
	}
	return "http"
}

// Port returns the TCP port for the running server.
// This is useful in tests where we allocate a random port by using ":0".
func (s *HTTPServer) Port() int {
	if s.ln == nil {
		return 0
	}
	return s.ln.Addr().(*net.TCPAddr).Port
}

// URL returns the local base URL of the running server.
func (s *HTTPServer) URL() string {
	scheme, port := s.Scheme(), s.Port()

	// Use localhost unless a domain is specified.
	domain := "localhost"
	if s.Domain != "" {
		domain = s.Domain
	}

	// Return without port if using standard ports.
	if (scheme == "http" && port == 80) || (scheme == "https" && port == 443) {
		return fmt.Sprintf("%s://%s", s.Scheme(), domain)
	}
	return fmt.Sprintf("%s://%s:%d", s.Scheme(), domain, s.Port())
}

// Open validates the server options and begins listening on the bind address.
func (s *HTTPServer) Open() (err error) {

	// Open a listener on our bind address.
	if s.Domain != "" {
		s.ln = autocert.NewListener(s.Domain)
	} else {
		if s.ln, err = net.Listen("tcp", s.Addr); err != nil {
			return err
		}
	}

	// Begin serving requests on the listener. We use Serve() instead of
	// ListenAndServe() because it allows us to check for listen errors (such
	// as trying to use an already open port) synchronously.
	go func() {
		err := s.server.Serve(s.ln)
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	return nil
}

// Close gracefully shuts down the server.
func (s *HTTPServer) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancel()
	return s.server.Shutdown(ctx)
}

// ListenAndServeTLSRedirect runs an HTTP server on port 80 to redirect users
// to the TLS-enabled port 443 server.
func ListenAndServeTLSRedirect(domain string) error {
	server := &http.Server{
		Addr:         ":80",
		Handler:      http.HandlerFunc(redirectTLS(domain)),
		ReadTimeout:  ReadTimeout,
		WriteTimeout: WriteTimeout,
		IdleTimeout:  IdleTimeout,
	}

	return server.ListenAndServe()
}

func redirectTLS(domain string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://"+domain, http.StatusFound)
	}
}

// ListenAndServeDebug runs an HTTP server with /debug endpoints (e.g. pprof, vars).
func ListenAndServeDebug() error {
	mux := http.NewServeMux()

	server := &http.Server{
		Addr:         ":6060",
		Handler:      mux,
		ReadTimeout:  ReadTimeout,
		WriteTimeout: WriteTimeout,
		IdleTimeout:  IdleTimeout,
	}

	return server.ListenAndServe()
}
