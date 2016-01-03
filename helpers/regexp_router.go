package helpers

import (
  "net/http"
  "regexp"
)

type routeMap struct {
  pattern *regexp.Regexp
  handler http.Handler
}

type RegexpRouter struct {
  routes []*routeMap
}

// func (h *RegexpRouter) Handler(pattern *regexp.Regexp, handler http.Handler) {
//   h.routes = append(h.routes, &routeMap{pattern, handler})
// }

func (h *RegexpRouter) AddRoute(pattern string, handler func(http.ResponseWriter, *http.Request)) {
  h.routes = append(h.routes, &routeMap{regexp.MustCompile(pattern), http.HandlerFunc(handler)})
}

func (h *RegexpRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  for _, route := range h.routes {
    if route.pattern.MatchString(r.URL.Path) {
      route.handler.ServeHTTP(w, r)
      return
    }
  }

  http.NotFound(w, r)
}
