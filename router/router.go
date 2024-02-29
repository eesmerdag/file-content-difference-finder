package router

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"

	"file-diff-finder/cache"
	fdf "file-diff-finder/diff_finder"
)

var (
	httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_request_duration_seconds",
		Help: "Duration of HTTP requests.",
	}, []string{"path"})
)

type ErrorResp struct {
	Message string
	Code    int
}

type Router struct {
	logger *logrus.Logger
	router *mux.Router
	cache  cache.ICache
	srv    fdf.IFileDiffFinder
}

func NewRouter(logger *logrus.Logger, cache cache.ICache, srv fdf.IFileDiffFinder) *Router {
	router := mux.NewRouter()

	r := &Router{
		logger: logger,
		router: router,
		cache:  cache,
		srv:    srv,
	}

	router.HandleFunc("/diff", r.Diff).Methods(http.MethodGet)
	router.HandleFunc("/live", r.Liveness).Methods(http.MethodGet)
	router.HandleFunc("/ready", r.Readiness).Methods(http.MethodGet)
	router.Handle("/metrics", promhttp.Handler())
	router.Use(panicRecovery)
	router.Use(prometheusMiddleware)

	return r
}

func panicRecovery(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				errorResponse(w, "unexpected internal error", http.StatusInternalServerError)
				return
			}
		}()

		h.ServeHTTP(w, r)
	})
}

func prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()
		timer := prometheus.NewTimer(httpDuration.WithLabelValues(path))
		next.ServeHTTP(w, r)
		timer.ObserveDuration()
	})
}

func (rt *Router) Diff(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var t RequestPayload
	err := decoder.Decode(&t)
	if err != nil {
		errorResponse(w, "error on parsing payload", http.StatusBadRequest)
		return
	}

	if err = t.Valid(); err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = rt.srv.ValidateVersion(t.Version)
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, updatingAnotherRequest := rt.cache.Get(strconv.Itoa(t.Version))
	if updatingAnotherRequest {
		errorResponse(w, "another process is getting the diff", http.StatusUnprocessableEntity)
		return
	}

	rt.cache.Set(strconv.Itoa(t.Version), true)
	defer rt.cache.Delete(strconv.Itoa(t.Version))

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	delta, err := rt.srv.Diff(ctx, *t.Text)
	if err != nil {
		rt.logger.Error("error on running diff", err.Error())
		errorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := Response{
		CurrentVersion: rt.srv.Version(),
		UpdatedVersion: t.Version,
		Delta:          delta,
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		rt.logger.Error("error on answering diff response", err.Error())
		errorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (rt *Router) Liveness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func (rt *Router) Readiness(w http.ResponseWriter, r *http.Request) {
	_, exists := rt.cache.Get(cache.ReadinessKey)
	if !exists {
		rt.logger.Warn("redis has no readiness key check, service is restarting")
		errorResponse(w, "error cache is not ok", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func (rt *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	rt.router.ServeHTTP(w, req)
}

func errorResponse(w http.ResponseWriter, message string, code int) {
	errObj := ErrorResp{
		Message: message,
		Code:    code,
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(errObj)
}
