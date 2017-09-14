package controller

import (
	"app/model"
	"app/shared/passhash"
	"fmt"
	"net/http"

	"encoding/json"

	"regexp"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// CreateUser creating a user
func CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var err error

	responsevalue := make(map[string]interface{})

	// get visitorid
	visitorid := r.Header.Get("visitorid")
	if visitorid != "guest" {
		responsevalue["error"] = "Forbidden, already have account. Should log out first"
		uj, _ := json.Marshal(responsevalue)
		w.WriteHeader(http.StatusForbidden)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "%s", uj)
		return
	}

	// Stub an user to be populated from the body
	u := model.Form{}

	// Populate the user data
	json.NewDecoder(r.Body).Decode(&u)
	u.UserName = strings.ToLower(u.UserName)
	u.Email = strings.ToLower(u.Email)

	///////////////////////// Validation /////////////////////////////////////
	if u.FirstName == "" {
		responsevalue["first_name"] = "This field is required"
	} else if matched, _ := regexp.MatchString("^[A-Za-z]([A-Za-z]{0,16}[A-Za-z])$", u.FirstName); !matched {
		responsevalue["first_name"] = "First Name only allowed 2-18 alphanumeric characters with underscore between"
	}
	if u.LastName == "" {
		responsevalue["last_name"] = "This field is required"
	} else if matched, _ := regexp.MatchString("^[A-Za-z]([A-Za-z]{0,16}[A-Za-z])$", u.LastName); !matched {
		responsevalue["last_name"] = "Last Name only allowed 2-18 alphanumeric characters with underscore between"
	}
	result, _ := model.GetUser(u.UserName)
	if u.UserName == "" {
		responsevalue["username"] = "This field is required"
	} else if matched, _ := regexp.MatchString("^[A-Za-z]([A-Za-z0-9_]{4,18}[A-Za-z0-9])$", u.UserName); !matched {
		responsevalue["username"] = "Username only allowed 6-20 lower case alphanumeric characters with underscore between"
	} else if result.UserName == u.UserName {
		responsevalue["username"] = "Username already taken, choose another username"
	}
	result, _ = model.GetUser(u.Email)
	if u.Email == "" {
		responsevalue["email"] = "This field is required"
	} else if result.Email == u.Email {
		responsevalue["email"] = "Email already taken, choose another email"
	}
	// ^[/\S+/]{8,32}$
	if u.Password == "" {
		responsevalue["password"] = "This field is required"
	} else if matched, _ := regexp.MatchString("^[/\\S+/]{8,32}$", u.Password); !matched {
		responsevalue["password"] = "Password only allowed 8-32 any non-whitespace characters"
	}
	if u.ConfirmPassword == "" {
		responsevalue["con_pass"] = "This field is required"
	} else if matched, _ := regexp.MatchString("^[/\\S+/]{8,32}$", u.ConfirmPassword); !matched {
		responsevalue["con_pass"] = "Password only allowed 8-32 any non-whitespace characters"
	} else if u.Password != u.ConfirmPassword {
		responsevalue["con_pass"] = "Confirm Password is not matched"
	}
	for _, value := range responsevalue {
		if value == "Username already taken, choose another username" || value == "Email already taken, choose another email" {
			w.WriteHeader(http.StatusConflict)
			uj, _ := json.Marshal(responsevalue)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, "%s", uj)
			return
		} else if value != "" {
			w.WriteHeader(http.StatusBadRequest)
			uj, _ := json.Marshal(responsevalue)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, "%s", uj)
			return
		}
	}

	////////////////////////// Creating User /////////////////////////////
	// Hashing pass
	u.Password, err = passhash.HashString(u.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Create User
	err = model.CreateUser(u.UserName, u.FirstName, u.LastName, u.Email, u.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%s", err)
		return
	}
	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(u)

	// Write content-type, status code, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s", uj)
}

// RetrieveUser getting a user by username
func RetrieveUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// get username
	username := p.ByName("username")

	// get visitorid
	visitorid := r.Header.Get("visitorid")

	// Lowering case for username, alphanumeric and underscore, 6 - 20 characters allowed
	matched, err := regexp.MatchString("^[A-Za-z]([A-Za-z0-9_]{4,18}[A-Za-z0-9])$", username)
	if !matched {
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else if err == nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "alphanumeric, underscore, and hypens")
			return
		}
	} else {
		username = strings.ToLower(username)
	}

	// Retrieve user from database
	result, err := model.GetUser(username)

	// Determine if user exists
	if err == model.ErrNoResult {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "no user exists")
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if visitorid != result.UserName {
		result.Password = ""
	}

	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(result)

	// Write content-type, status code, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", uj)
}

// UpdateUser editting a user by username. Not allowed changing username
func UpdateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var err error

	username := p.ByName("username")

	visitorid := r.Header.Get("visitorid")
	if visitorid != username {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "forbidden editting other account")
		return
	}

	// Stub an user to be populated from the body
	u := model.UpdateForm{}

	// Populate the user data
	json.NewDecoder(r.Body).Decode(&u)

	result, err := model.GetUser(username)
	u.UserName = result.UserName

	// Check firstname and only allowed 2 to 20 characters
	_, err = regexp.MatchString("^[A-Za-z]([A-Za-z]{0,18}[A-Za-z])$", u.NewFirstName)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "First name, 2 - 20 alphabet")
		return
	}

	// Check lastname and only allowed 2 to 20 characters
	_, err = regexp.MatchString("^[A-Za-z]([A-Za-z]{0,18}[A-Za-z])$", u.NewLastName)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Last name, 2 - 20 alphabet")
		return
	}

	// Check Old Password
	if u.OldPassword != result.Password {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Old password is invalid")
		return
	}

	// Check pass confirm
	if u.NewPassword != u.ConfirmPassword {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "password confirmation is invalid")
		return
	}
	// Hashing pass
	u.NewPassword, err = passhash.HashString(u.NewPassword)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Check Email
	if u.NewEmail != result.Email {
		result, err = model.GetUser(u.NewEmail)
		if result.Email == u.NewEmail {
			w.WriteHeader(http.StatusConflict)
			fmt.Fprintf(w, "email already existed")
			return
		}
	}

	// Create User
	err = model.UpdateUser(u.UserName, u.NewFirstName, u.NewLastName, u.NewEmail, u.NewPassword)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%s", err)
		return
	}
	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(u)

	// Write content-type, status code, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", uj)
}

/*
// DeleteUser deleting a user by username
func DeleteUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Delete user")
}
*/
