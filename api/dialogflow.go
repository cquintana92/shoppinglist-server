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
	ResponseId  string `json:"responseId"`
	QueryResult struct {
		QueryText  string `json:"queryText"`
		Parameters struct {
			Ingredient string `json:"Ingredient"`
		} `json:"parameters"`
		AllRequiredParamsPresent bool   `json:"allRequiredParamsPresent"`
		FulfillmentText          string `json:"fulfillmentText"`
		FulfillmentMessages      []struct {
			Text struct {
				Text []string `json:"text"`
			} `json:"text"`
		} `json:"fulfillmentMessages"`
		OutputContexts []struct {
			Name       string `json:"name"`
			Parameters struct {
				NoInput            float64 `json:"no-input"`
				NoMatch            float64 `json:"no-match"`
				Ingredient         string  `json:"Ingredient"`
				IngredientOriginal string  `json:"Ingredient.original"`
			} `json:"parameters"`
		} `json:"outputContexts"`
		Intent struct {
			Name           string `json:"name"`
			DisplayName    string `json:"displayName"`
			EndInteraction bool   `json:"endInteraction"`
		} `json:"intent"`
		IntentDetectionConfidence float64 `json:"intentDetectionConfidence"`
		LanguageCode              string  `json:"languageCode"`
		SentimentAnalysisResult   struct {
			QueryTextSentiment struct {
				Score     float64 `json:"score"`
				Magnitude float64 `json:"magnitude"`
			} `json:"queryTextSentiment"`
		} `json:"sentimentAnalysisResult"`
	} `json:"queryResult"`
	OriginalDetectIntentRequest struct {
		Source  string `json:"source"`
		Payload struct {
		} `json:"payload"`
	} `json:"originalDetectIntentRequest"`
	Session string `json:"session"`
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

func dialogFlowHandler(w http.ResponseWriter, r *http.Request) {
	log.Logger.Info("Received DialogFlow request")

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Logger.Errorf("Error reading DialogFlow request body: %+v", err)
		return
	}

	request := &DialogFlowRequest{}
	err = json.Unmarshal(body, request)
	if err != nil {
		log.Logger.Errorf("Error reading DialogFlow as JSON: %+v", err)
		return
	}

	log.Logger.Debugf("%+v", request)
	ingredient := request.QueryResult.Parameters.Ingredient
	log.Logger.Infof("DialogFlow ingredient: %s", ingredient)

	err = storage.New(ingredient)
	var dialogFlowResponse string
	if err != nil {
		if err == storage.ItemAlreadyExistsError {
			log.Logger.Warnf("Item [%s] already exists", ingredient)
			dialogFlowResponse = fmt.Sprintf("No he añadido %s porque ya estaba", ingredient)
		} else {
			log.Logger.Errorf("Error creating a new record: %+v", err)
			dialogFlowResponse = fmt.Sprintf("No he podido añadir %s porque ha habido un error. Carlos, te toca mirarlo", ingredient)
		}
	} else {
		dialogFlowResponse = fmt.Sprintf("Acabo de añadir %s a tu lista de la compra", ingredient)
	}

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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(asBytes)
}
