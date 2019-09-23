package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"fullstack/api/auth"
	"fullstack/api/models"
	"fullstack/responses"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"

	formaterror "fullstack/api/utils"
)

func (server Server) CreatePost(w http.ResponseWriter, r *http.Request) {
	//read value from request
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}

	post := models.Post{}

	//convert body to post structure
	err = json.Unmarshal(body, &post)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	//post validation
	post.Prepare()
	err = post.Validate()

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	//extract token ID from request
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	//validate user id
	if uid != post.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	postCreated, err := post.SavePost(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	//save model
	w.Header().Set("Lacation", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, postCreated.ID))
	responses.JSON(w, http.StatusCreated, postCreated)

}

func (server Server) GetPosts(w http.ResponseWriter, r *http.Request) {

	post := models.Post{}

	posts, error := post.FindAllPosts(server.DB)

	if error != nil {
		responses.ERROR(w, http.StatusInternalServerError, error)
		return
	}

	responses.JSON(w, http.StatusOK, posts)
}

func (server Server) GetPost(w http.ResponseWriter, r *http.Request) {

	//get value from request
	vars := mux.Vars(r)

	pid, err := strconv.ParseUint(vars["id"], 10, 64)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	post := models.Post{}

	postReceived, err := post.FindPostByID(server.DB, pid)

	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, postReceived)
}

func (server Server) UpdatePost(w http.ResponseWriter, r *http.Request) {

}

func (server *Server) DeletePost(w http.ResponseWriter, r *http.Request) {

}
