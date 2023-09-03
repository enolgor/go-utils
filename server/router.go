package server

import (
	"fmt"
	"net/http"

	"github.com/enolgor/go-utils/server/path"
)

type Router struct {
	routes          []route
	notFoundHandler http.HandlerFunc
	recoverHandler  http.HandlerFunc
	preFilters      []http.HandlerFunc
	postFilters     []http.HandlerFunc
}

type route struct {
	method   string
	pathExpr string
	matcher  func(string, map[any]string) bool
	handler  http.HandlerFunc
}

func NewRouter() *Router {
	return &Router{
		routes:          []route{},
		notFoundHandler: defaultNotFoundHandler,
		recoverHandler:  defaultRecoverHandler,
		preFilters:      []http.HandlerFunc{},
		postFilters:     []http.HandlerFunc{},
	}
}

func (r *Router) register(method string, pathExpr string, handler http.HandlerFunc) *Router {
	matcher, err := path.Matcher(pathExpr)
	if err != nil {
		panic(err)
	}
	r.routes = append(r.routes, route{method, pathExpr, matcher, handler})
	return r
}

func (r *Router) Get(pathExpr string, handler http.HandlerFunc) *Router {
	return r.register("GET", pathExpr, handler)
}

func (r *Router) Post(pathExpr string, handler http.HandlerFunc) *Router {
	return r.register("POST", pathExpr, handler)
}

func (r *Router) Put(pathExpr string, handler http.HandlerFunc) *Router {
	return r.register("PUT", pathExpr, handler)
}

func (r *Router) Patch(pathExpr string, handler http.HandlerFunc) *Router {
	return r.register("PATCH", pathExpr, handler)
}

func (r *Router) Delete(pathExpr string, handler http.HandlerFunc) *Router {
	return r.register("DELETE", pathExpr, handler)
}

func (r *Router) NotFoundHandler(handler http.HandlerFunc) *Router {
	r.notFoundHandler = handler
	return r
}

func (r *Router) RecoverHandler(handler http.HandlerFunc) *Router {
	r.recoverHandler = handler
	return r
}

func (r *Router) PreFilters(filters ...http.HandlerFunc) *Router {
	r.preFilters = filters
	return r
}

func (r *Router) PostFilters(filters ...http.HandlerFunc) *Router {
	r.postFilters = filters
	return r
}

func (r *Router) SubRoute(pathExpr string, router *Router) *Router {
	var err error
	for i := range router.routes {
		router.routes[i].pathExpr = pathExpr + router.routes[i].pathExpr
		if router.routes[i].matcher, err = path.Matcher(router.routes[i].pathExpr); err != nil {
			panic(err)
		}
	}
	r.register("ANY", pathExpr+"(.*)", router.ServeHTTP)
	return r
}

var defaultNotFoundHandler = func(w http.ResponseWriter, req *http.Request) {
	Response(w).Status(http.StatusNotFound).WithBody(fmt.Sprintf("%s %s not found", req.Method, req.URL.Path)).AsTextPlain()
}

var defaultRecoverHandler = func(w http.ResponseWriter, req *http.Request) {
	err := Recover(req)
	if err == nil {
		err = "uknown"
	}
	Response(w).Status(http.StatusInternalServerError).WithBody(fmt.Sprintf("internal server error: %s", err)).AsTextPlain()
}

type routerContextKey int

const (
	pathParamsKey routerContextKey = iota
	panicKey
)

func PathParams(req *http.Request) map[any]string {
	values := map[any]string{}
	GetContextValue(req, pathParamsKey, &values)
	return values
}

func Recover(req *http.Request) any {
	var err any
	GetContextValue(req, panicKey, &err)
	return err
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			AddContextValue(req, panicKey, rec)
			r.recoverHandler(w, req)
			r.applyPostFilters(w, req)
		}
	}()
	pathParams := make(map[any]string)
	var handler http.HandlerFunc
	for _, route := range r.routes {
		if (route.method == "ANY" || req.Method == route.method) && route.matcher(req.URL.Path, pathParams) {
			AddContextValue(req, pathParamsKey, pathParams)
			handler = route.handler
			break
		}
	}
	if handler == nil {
		handler = r.notFoundHandler
	}
	rw := NewResponseWriter(w)
	r.applyPreFilters(rw, req)
	handler(rw, req)
	r.applyPostFilters(rw, req)
}

func (r *Router) applyPreFilters(w http.ResponseWriter, req *http.Request) {
	for _, filter := range r.preFilters {
		filter(w, req)
	}
}

func (r *Router) applyPostFilters(w http.ResponseWriter, req *http.Request) {
	for _, filter := range r.postFilters {
		filter(w, req)
	}
}
