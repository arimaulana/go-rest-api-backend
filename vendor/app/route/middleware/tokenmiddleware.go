package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/julienschmidt/httprouter"

	"app/controller"
)

// AnonUser : validating the token and sub
func AnonUser(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor, func(token *jwt.Token) (interface{}, error) {
			return controller.MySigningKey, nil
		})

		responsevalue := make(map[string]interface{})

		if err == nil {
			if token.Valid {
				claims, ok := token.Claims.(jwt.MapClaims)
				if !ok {
					responsevalue["token"] = "Internal error parsing token to claim"
					uj, _ := json.Marshal(responsevalue)
					w.WriteHeader(http.StatusInternalServerError)
					w.Header().Set("Content-Type", "application/json")
					fmt.Fprintf(w, "%s", uj)
				}
				visitorid := claims["sub"].(string)
				r.Header.Set("visitorid", visitorid)
				next(w, r, ps)
			} else if !token.Valid {
				responsevalue["token"] = "Token is not valid, please log in first"
				uj, _ := json.Marshal(responsevalue)
				w.WriteHeader(http.StatusUnauthorized)
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintf(w, "%s", uj)
			}
		} else if err == request.ErrNoTokenInRequest {
			r.Header.Set("visitorid", "guest")
			next(w, r, ps)
		} else {
			responsevalue["token"] = "Unauthorized access to this resource"
			uj, _ := json.Marshal(responsevalue)
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, "%s", uj)
		}
	}
}

// AuthUser : validating the token and sub
func AuthUser(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor, func(token *jwt.Token) (interface{}, error) {
			return controller.MySigningKey, nil
		})

		responsevalue := make(map[string]interface{})

		if err == nil {
			if token.Valid {
				claims, ok := token.Claims.(jwt.MapClaims)
				if !ok {
					responsevalue["token"] = "Internal error parsing token to claim"
					uj, _ := json.Marshal(responsevalue)
					w.WriteHeader(http.StatusInternalServerError)
					w.Header().Set("Content-Type", "application/json")
					fmt.Fprintf(w, "%s", uj)
				}
				visitorid := claims["sub"].(string)
				r.Header.Set("visitorid", visitorid)
				next(w, r, ps)
			} else if !token.Valid {
				responsevalue["token"] = "Token is not valid"
				uj, _ := json.Marshal(responsevalue)
				w.WriteHeader(http.StatusUnauthorized)
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintf(w, "%s", uj)
			}
		} else {
			responsevalue["token"] = "Unauthorized access to this resource"
			uj, _ := json.Marshal(responsevalue)
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, "%s", uj)
		}
	}
}
