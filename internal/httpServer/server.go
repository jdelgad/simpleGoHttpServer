package httpServer

import (
	"context"
	"net/http"
	"regexp"
)

type route struct {
	pattern *regexp.Regexp
	handler http.HandlerFunc
}

// Server represents a webserver responsible for dispatching incoming
// requests to designated routes
type Server struct {
	addr string
	routes []route
	server http.Server
}

func NewServer(addr string) *Server {
	return &Server{
		server:  http.Server{
			Addr: addr,
		},
	}
}

// Handle registers routes in order as they are added
func (s *Server) Handle(pattern string, handler http.HandlerFunc) {
	re := regexp.MustCompile(pattern)
	route := route{
		pattern: re,
		handler: handler,
	}
	s.routes = append(s.routes, route)
}

func (s *Server) ListenAndServe() error {
	s.server.Handler = s
	return s.server.ListenAndServe()
}

// ServeHTTP dispatches handlers based on the matching route, if
// there are numerous matches, the first route and handler registered
// is dispatched
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, re := range s.routes {
		if matches := re.pattern.MatchString(r.URL.Path); matches {
			re.handler(w, r)
			return
		}
	}
	w.WriteHeader(http.StatusBadRequest)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}