package api

import (
	"encoding/json"
	"io"
	"net/http"
	"shoppinglistserver/log"
	"shoppinglistserver/storage"
	"shoppinglistserver/utils"
)

const todoistAppIdHeader = "X-App-ID"

type TodoistConfig struct {
	Enabled  bool
	AppId    string
	Endpoint string
}

type todoistRequest struct {
	EventName string                   `json:"event_name"`
	EventData todoistCreateItemRequest `json:"event_data"`
}

type todoistCreateItemRequest struct {
	Content string `json:"content"`
}

func todoist(config TodoistConfig) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		log.Logger.Info("Received TODOIST event request")

		if !validateTodoistToken(config.AppId, r) {
			log.Logger.Errorf("AppID does not match")
			http.Error(w, "AppID does not match", 400)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Logger.Errorf("Could not read the request body: %+v", err)
			http.Error(w, "Could not create a record", 500)
			return
		}

		req := todoistRequest{}
		err = json.Unmarshal(body, &req)
		log.Logger.Debugf("TodoistRequest eventName: %s", req.EventName)
		if err != nil {
			log.Logger.Errorf("Could not convert to JSON: %+v", err)
			http.Error(w, "Could not create a record", 500)
			return
		}

		name := req.EventData.Content
		if name == "" {
			log.Logger.Errorf("Recieved a request to create with an empty name")
			http.Error(w, "Could not create a record", 400)
			return
		}
		name = utils.SanitizeName(name)
		log.Logger.Debugf("Received %+v", req)
		log.Logger.Infof("Creating item: %s", name)
		err = storage.New(name)
		if err != nil {
			if err == storage.ItemAlreadyExistsError {
				log.Logger.Errorf("Item [%s] already exists", name)
				http.Error(w, "Item already exists", 409)
			} else {
				log.Logger.Errorf("Error creating a new record: %+v", err)
				http.Error(w, "Could not create a record", 500)
			}
		} else {
			respondAll(w, r)
		}
	}
}

func validateTodoistToken(expected string, r *http.Request) bool {
	appIdHeader := r.Header.Get(todoistAppIdHeader)
	if appIdHeader == expected {
		return true
	} else {
		log.Logger.Warnf("Todoist token did not match [expected=%s=] [got=%s]", expected, appIdHeader)
		return false
	}
}
