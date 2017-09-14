package controller

import (
	"app/model"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// GetPostByID : gets post by post id
func GetPostByID(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	postid := p.ByName("id")

	post, err := model.GetPostByID(postid)
	// Determine if user exists
	if err == model.ErrNoResult {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "post doesnt exists")
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	uj, _ := json.Marshal(post)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", uj)
}

// GetPostByFeed : gets post for a user
func GetPostByFeed(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	feedid := p.ByName("feed")

	posts, err := model.GetPostByFeed(feedid)
	if err != nil {
		log.Println(err)
		posts = []model.Post{}
	}

	uj, _ := json.Marshal(posts)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", uj)
}

// CreatePostFeed creating post from user page
func CreatePostFeed(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var err error

	// Stub an user to be populated from the body
	u := model.Post{}

	// Populate the user data
	json.NewDecoder(r.Body).Decode(&u)

	if len(u.PostContent) == 0 {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "We cant process blank post")
		return
	}

	feedid := p.ByName("feed")
	visitorid := r.Header.Get("visitorid")
	content := u.PostContent

	err = model.CreatePostFeed(feedid, visitorid, content)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(u)

	// Write content-type, status code, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s", uj)
}

// UpdatePost : updates a post
func UpdatePost(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var err error

	// Stub an user to be populated from the body
	u := model.Post{}

	// Populate the user data
	json.NewDecoder(r.Body).Decode(&u)

	if len(u.PostContent) == 0 {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "We cant process blank post")
		return
	}

	postid := p.ByName("id")
	visitorid := r.Header.Get("visitorid")
	content := u.PostContent

	err = model.UpdatePost(content, visitorid, postid)
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

// DeletePost : deletes a post
func DeletePost(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	postid := p.ByName("id")

	// later, use visitorid from token jwt
	visitorid := r.Header.Get("visitorid")

	err := model.DeletePost(visitorid, postid)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if err == nil {
		w.WriteHeader(200)
		http.Redirect(w, r, "/", 200)
	}
}
