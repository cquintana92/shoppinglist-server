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

type ApiConfig struct {
	Port           int
	SecretEndpoint string
	SecretBearer   string
	TodoistConfig  TodoistConfig
}

func Run(config *ApiConfig) error {

	if config.SecretBearer != "" {
		log.Logger.Info("Secret bearer set")
		secretBearerAuthorization = config.SecretBearer
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

	if config.SecretEndpoint != "" {
		log.Logger.Info("Secret endpoint set")
		endpoint := fmt.Sprintf("/%s", config.SecretEndpoint)
		r.HandleFunc(endpoint, create).Methods(http.MethodPost)
	}

	if config.TodoistConfig.Enabled {
		log.Logger.Info("Todoist endpoint set")
		endpoint := fmt.Sprintf("/%s", config.TodoistConfig.Endpoint)
		r.HandleFunc(endpoint, todoist(config.TodoistConfig)).Methods(http.MethodPost)
	}

	addr := fmt.Sprintf("0.0.0.0:%d", config.Port)
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
