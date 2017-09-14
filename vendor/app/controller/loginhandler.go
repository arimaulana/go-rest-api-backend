package controller

import (
	"log"
	"net/http"
	"regexp"

	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"

	"app/model"
	"app/shared/passhash"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// MySigningKey sign key
var MySigningKey = []byte("RahasiaBro")

// LoginHandler : assigning token after login
func LoginHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var user model.Form

	responsevalue := make(map[string]interface{})

	// get visitorid and block if already have token
	visitorid := r.Header.Get("visitorid")
	log.Println(visitorid)
	if visitorid != "guest" {
		responsevalue["error"] = "Forbidden, already logged in. Should log out first"
		uj, _ := json.Marshal(responsevalue)
		w.WriteHeader(http.StatusForbidden)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s", uj)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		responsevalue["error"] = "Error in decoding request body"
		uj, _ := json.Marshal(responsevalue)
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s", uj)
		return
	}

	username := strings.ToLower(user.UserName)

	/////////////////////////// Validation ////////////////////////////////////////
	result, err := model.GetUser(username)
	if username == "" {
		responsevalue["username"] = "This field is required"
	} else if matched, _ := regexp.MatchString("^[A-Za-z]([A-Za-z0-9_]{4,18}[A-Za-z0-9])$", username); !matched {
		responsevalue["username"] = "Username only allowed 6-20 lower case alphanumeric characters with underscore between"
	} else if err == model.ErrNoResult {
		responsevalue["username"] = "No user exist"
	}
	if user.Password == "" {
		responsevalue["password"] = "This field is required"
	} else if !passhash.MatchString(result.Password, user.Password) {
		responsevalue["password"] = "Invalid Credentials"
	}
	if result.StatusID != 1 && passhash.MatchString(result.Password, user.Password) {
		responsevalue["username"] = "User is inactive"
	}
	for _, value := range responsevalue {
		if value != "" {
			w.WriteHeader(http.StatusBadRequest)
			uj, _ := json.Marshal(responsevalue)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, "%s", uj)
			return
		}
	}

	if result.UserName == username && passhash.MatchString(result.Password, user.Password) {
		if result.StatusID == 1 {
			token := jwt.New(jwt.SigningMethodHS256)
			claims := make(jwt.MapClaims)
			claims["exp"] = time.Now().Add(time.Minute * 30).Unix()
			claims["iat"] = time.Now().Unix()
			claims["sub"] = result.UserName
			token.Claims = claims

			tokenString, err := token.SignedString(MySigningKey)
			if err != nil {
				responsevalue["error"] = "Error while signing the token"
				uj, _ := json.Marshal(responsevalue)
				w.WriteHeader(http.StatusInternalServerError)
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintf(w, "%s", uj)
				return
			} else if err == nil {
				responsevalue["token"] = tokenString
				responsevalue["user"] = result.UserName
				uj, _ := json.Marshal(responsevalue)
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintf(w, "%s", uj)
				return
			}
		}
	}
}
