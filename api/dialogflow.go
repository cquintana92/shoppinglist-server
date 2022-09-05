package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"shoppinglistserver/log"
	"shoppinglistserver/storage"
)

type DialogFlowRequest struct {
	QueryResult struct {
		QueryText  string `json:"queryText"`
		Parameters struct {
			Ingredient string `json:"Ingredient"`
		} `json:"parameters"`
		Intent struct {
			Name           string `json:"name"`
			DisplayName    string `json:"displayName"`
			EndInteraction bool   `json:"endInteraction"`
		} `json:"intent"`
	} `json:"queryResult"`
}

type DialogFlowResponse struct {
	FulfilmentMessages []DialogFlowFulfilmentMessage `json:"fulfillmentMessages"`
	Payload            DialogFlowResponsePayload     `json:"payload"`
}

type DialogFlowFulfilmentMessage struct {
	Text DialogFlowFulfilmentTextMessage `json:"text"`
}

type DialogFlowFulfilmentTextMessage struct {
	Text []string `json:"text"`
}

type DialogFlowResponsePayload struct {
	Google DialogFlowGoogleResponsePayload `json:"google"`
}

type DialogFlowGoogleResponsePayload struct {
	ExpectUserResponse bool                           `json:"expectUserResponse"`
	RichResponse       DialogFlowResponseRichResponse `json:"richResponse"`
}

type DialogFlowResponseRichResponse struct {
	Items []DialogFlowResponseItem `json:"items"`
}

type DialogFlowResponseItem struct {
	SimpleResponse DialogFlowResponseSimpleResponse `json:"simpleResponse"`
}

type DialogFlowResponseSimpleResponse struct {
	TextToSpeech string `json:"textToSpeech"`
}

func processDialogFlowRequest(body []byte) string {
	request := &DialogFlowRequest{}
	err := json.Unmarshal(body, request)
	if err != nil {
		log.Logger.Errorf("Error reading DialogFlow as JSON: %+v", err)
		return "No he podido hacer lo que me has pedido porque ha habido un error"
	}

	log.Logger.Debugf("%+v", request)
	ingredient := request.QueryResult.Parameters.Ingredient
	log.Logger.Infof("DialogFlow ingredient: %s", ingredient)

	if ingredient == "" {
		return "Perdona, no he entendido qué es lo que quieres que apunte. Igual ha habido un error"
	}

	err = storage.New(ingredient)
	if err != nil {
		if err == storage.ItemAlreadyExistsError {
			log.Logger.Warnf("Item [%s] already exists", ingredient)
			return fmt.Sprintf("No he añadido %s porque ya estaba", ingredient)
		} else {
			log.Logger.Errorf("Error creating a new record: %+v", err)
			return fmt.Sprintf("No he podido añadir %s porque ha habido un error.", ingredient)
		}
	} else {
		return fmt.Sprintf("Acabo de añadir %s a tu lista de la compra", ingredient)
	}
}

func dialogFlowHandler(w http.ResponseWriter, r *http.Request) {
	log.Logger.Info("Received DialogFlow request")

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Logger.Errorf("Error reading DialogFlow request body: %+v", err)
		return
	}

	// We must always return a response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	dialogFlowResponse := processDialogFlowRequest(body)

	res := &DialogFlowResponse{
		FulfilmentMessages: []DialogFlowFulfilmentMessage{
			{Text: DialogFlowFulfilmentTextMessage{
				Text: []string{dialogFlowResponse},
			}},
		},
		Payload: DialogFlowResponsePayload{
			Google: DialogFlowGoogleResponsePayload{
				ExpectUserResponse: false,
				RichResponse: DialogFlowResponseRichResponse{
					Items: []DialogFlowResponseItem{
						{SimpleResponse: DialogFlowResponseSimpleResponse{
							TextToSpeech: dialogFlowResponse,
						}},
					},
				},
			},
		},
	}

	asBytes, err := json.Marshal(res)
	if err != nil {
		log.Logger.Errorf("Error converting DialogFlow response to JSON: %+v", err)
		return
	}

	w.Write(asBytes)
}
