package controllers

import (
	"encoding/json"
	"fullstack/api/auth"
	"fullstack/api/models"
	"fullstack/responses"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"net/http"
)

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	user := models.User{}

	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
}

func (s *Server) SignIn(email, password string) (string, error) {
	var err error

	user := models.User{}

	err = s.DB.Debug().Model(models.User{}).Where("email = ?", email).Take(&user).Error

	if err != nil {
		return "", nil
	}

	err = models.VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
	return auth.CreateToken(user.ID)
}
