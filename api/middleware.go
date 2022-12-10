package api

import (
	"net/http"
	"shoppinglistserver/log"
	"strings"
)

type handler func(w http.ResponseWriter, r *http.Request)

func isBearerCorrect(r *http.Request) error {
	authorization := r.Header.Get("Authorization")
	if authorization == "" {
		return noAuthorizationError
	}

	if !strings.Contains(authorization, "Bearer ") {
		return badAuthorizationFormatError
	}

	withoutBearer := strings.Replace(authorization, "Bearer ", "", -1)
	if withoutBearer == secretBearerAuthorization {
		return nil
	} else {
		return invalidAuthorizationError
	}
}

func withBearer(h handler) handler {
	return func(w http.ResponseWriter, r *http.Request) {
		if secretBearerAuthorization != "" {
			if err := isBearerCorrect(r); err != nil {
				log.Logger.Errorf("WithBearer failed: %+v", err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}
		h(w, r)
	}
}
