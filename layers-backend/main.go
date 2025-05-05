package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"layersapi/controllers"
	"layersapi/repositories/files/csv"
	"layersapi/services"
	"net/http"
)

func main() {
	r := mux.NewRouter()

	userRepository := csv.NewUserRepository()
	userService := services.NewUserService(userRepository)
	userController := controllers.NewUserController(*userService)

	r.HandleFunc("/users", userController.GetAllUsersHandler).Methods(http.MethodGet)
	r.HandleFunc("/users", userController.CreateUserHandler).Methods(http.MethodPost)
	r.HandleFunc("/findbyid", userController.GetUserByIdHandler).Methods(http.MethodGet)
	r.HandleFunc("/updatebyid/", userController.UpdateUserHandler).Methods(http.MethodPut)

	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", r)

}
