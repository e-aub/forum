package main

import (
	"fmt"
	"net/http"
	"sync"
)

// RouteHandler stores handlers for different HTTP methods
type RouteHandler struct {
	methods map[string]http.HandlerFunc
	mu      sync.RWMutex
}

// AddHandler adds a handler for a specific HTTP method
func (rh *RouteHandler) AddHandler(method string, handler http.HandlerFunc) {
	rh.mu.Lock()
	defer rh.mu.Unlock()
	if rh.methods == nil {
		rh.methods = make(map[string]http.HandlerFunc)
	}
	rh.methods[method] = handler
}

// ServeHTTP implements the http.Handler interface
func (rh *RouteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rh.mu.RLock()
	defer rh.mu.RUnlock()
	if handler, exists := rh.methods[r.Method]; exists {
		handler(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Router stores routes and their handlers
type Router struct {
	routes map[string]*RouteHandler
	mu     sync.RWMutex
}

// NewRouter creates a new Router
func NewRouter() *Router {
	return &Router{
		routes: make(map[string]*RouteHandler),
	}
}

// Handle registers a route with a specific method and handler
func (r *Router) Handle(path, method string, handler http.HandlerFunc) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.routes[path]; !exists {
		r.routes[path] = &RouteHandler{}
	}
	r.routes[path].AddHandler(method, handler)
}

// ServeHTTP dispatches the request to the appropriate route and handler
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if routeHandler, exists := r.routes[req.URL.Path]; exists {
		routeHandler.ServeHTTP(w, req)
	} else {
		http.NotFound(w, req)
	}
}

// Example Handlers
func getHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "GET request handled")
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "POST request handled")
}

func main() {
	router := NewRouter()

	// Register routes
	router.Handle("/route", http.MethodGet, getHandler)
	router.Handle("/route", http.MethodPost, postHandler)

	// Start server
	http.ListenAndServe(":8080", router)
}
