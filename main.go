package main

import (
	"fmt"
	"gocrud/servidor"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	// Create - Post
	// Read - Get
	// Update - Put
	// Delete - Delete

	router := mux.NewRouter()
	router.HandleFunc("/usuario", servidor.CriarUsuario).Methods(http.MethodPost)
	router.HandleFunc("/usuarios", servidor.RetornaUsuarios).Methods(http.MethodGet)
	router.HandleFunc("/usuario/{id}", servidor.RetornaUsuario).Methods(http.MethodGet)
	router.HandleFunc("/usuario/{id}", servidor.AtualizarUsuario).Methods(http.MethodPut)
	router.HandleFunc("/usuario/{id}", servidor.DeletarUsuario).Methods(http.MethodDelete)

	fmt.Println("Escutando na porta 500")
	log.Fatal(http.ListenAndServe(":500", router))
}
