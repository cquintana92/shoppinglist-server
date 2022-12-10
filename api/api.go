package api

import (
	"fmt"
	"net/http"
	"shoppinglistserver/log"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var secretBearerAuthorization = ""

func Run(port int, secretEndpoint string, secretBearer string) error {

	if secretBearer != "" {
		log.Logger.Info("Secret bearer set")
		secretBearerAuthorization = secretBearer
	} else {
		log.Logger.Warn("Starting API without secret bearer")
	}

	r := mux.NewRouter()
	r.HandleFunc("/health", health).Methods(http.MethodGet)
	r.HandleFunc("/", withBearer(getAll)).Methods(http.MethodGet)
	r.HandleFunc("/", withBearer(create)).Methods(http.MethodPost)
	r.HandleFunc("/", withBearer(deleteAllChecked)).Methods(http.MethodDelete)
	r.HandleFunc("/{id}", withBearer(toggleChecked)).Methods(http.MethodPatch)
	r.HandleFunc("/{id}", withBearer(update)).Methods(http.MethodPut)
	r.HandleFunc("/{id}", withBearer(deleteOne)).Methods(http.MethodDelete)
	r.HandleFunc("/{id}/position", withBearer(updatePosition)).Methods(http.MethodPut)

	if secretEndpoint != "" {
		log.Logger.Info("Secret endpoint set")
		endpoint := fmt.Sprintf("/%s", secretEndpoint)
		r.HandleFunc(endpoint, create).Methods(http.MethodPost)
	}

	addr := fmt.Sprintf("0.0.0.0:%d", port)
	srv := &http.Server{
		Handler: cors.AllowAll().Handler(r),
		Addr:    addr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Logger.Infof("[API] Started server at: %s", addr)
	return srv.ListenAndServe()
}
