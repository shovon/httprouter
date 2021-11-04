package httprouter

import (
	"fmt"
	"net/http"
	"strings"
)

type Router struct {
	methods        map[string]map[string]http.Handler
	paths          map[string]http.Handler
	delegates      map[string]http.Handler
	defaultHandler DefaultHandler
}

var _ http.Handler = Router{}

func (receiver Router) handleMethods(w http.ResponseWriter, r *http.Request) {
	method, ok := receiver.methods[r.Method]
	if ok {
		path, ok := method[r.URL.Path]
		if !ok {
			receiver.handlePaths(w, r)
		} else {
			path.ServeHTTP(w, r)
		}
	} else {
		receiver.handlePaths(w, r)
	}
}

func (receiver Router) handlePaths(w http.ResponseWriter, r *http.Request) {
	path, ok := receiver.paths[r.URL.Path]
	if ok {
		path.ServeHTTP(w, r)
	} else {
		receiver.handleDelegates(w, r)
	}
}

func (receiver Router) handleDelegates(w http.ResponseWriter, r *http.Request) {
	foundPath := false
	for path, handler := range receiver.delegates {
		if strings.HasPrefix(r.URL.Path, path) {
			foundPath = true
			r.URL.Path = r.URL.Path[len(path):]
			handler.ServeHTTP(w, r)
			break
		}
	}
	if !foundPath {
		receiver.handleDefault(w, r)
	}
}

func (receiver Router) handleDefault(w http.ResponseWriter, r *http.Request) {
	receiver.defaultHandler.ServeHTTP(w, r)
}

func (receiver Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	receiver.handleMethods(w, r)
}

func (receiver *Router) Register(handler http.Handler) {
	receiver.defaultHandler = AlternativeHandler(handler)
}

func (receiver *Router) RegisterPath(handler http.Handler, path string) error {
	if receiver.paths == nil {
		receiver.paths = make(map[string]http.Handler)
	}
	receiver.paths[path] = handler
	return nil
}

func (receiver *Router) RegisterPathMethod(handler http.Handler, path string, methods ...string) error {
	if receiver.methods == nil {
		receiver.methods = make(map[string]map[string]http.Handler)
	}
	for _, m := range methods {
		method, ok := receiver.methods[m]
		if !ok {
			receiver.methods[m] = make(map[string]http.Handler)
			method1, ok := receiver.methods[m]
			if !ok {
				return fmt.Errorf("failed to register a handler for method %q, to path %q", method, path)
			}
			method = method1
		}
		method[path] = handler
	}
	return nil
}

func (receiver *Router) Delegate(handler http.Handler, path string) error {
	if receiver.delegates == nil {
		receiver.delegates = make(map[string]http.Handler)
	}
	receiver.delegates[path] = handler
	return nil
}
