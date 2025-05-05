package controllers

import (
	"encoding/json"
	"layersapi/entities/dto"
	"layersapi/services"
	"net/http"
)

type UserController struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) *UserController {
	return &UserController{userService: userService}
}

func (u UserController) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	resData, err := u.userService.GetAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	res, err := json.Marshal(resData)

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

// still working on it
func (u UserController) GetUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}
	id := r.FormValue("id")
	resData, err := u.userService.GetById(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	res, err := json.Marshal(resData)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func (u UserController) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}
	name := r.FormValue("name")
	email := r.FormValue("email")

	var err = u.userService.Create(dto.CreateUser{name, email})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return

	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

}

func (u UserController) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
	}

	id := r.FormValue("id")
	name := r.FormValue("name")
	email := r.FormValue("email")

	var err = u.userService.Update(id, dto.UpdateUser{name, email})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

}
