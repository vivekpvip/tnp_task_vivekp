package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Certificate struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
	Owner   string `json:"owner"`
}

var certificates []Certificate
var nextID int = 1

func GetCertificateByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	for _, certificate := range certificates {
		if certificate.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(certificate)
			return
		}
	}

	http.Error(w, "Certificate not found", http.StatusNotFound)
}

func CreateCertificate(w http.ResponseWriter, r *http.Request) {
	var newCertificate Certificate

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&newCertificate); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newCertificate.ID = nextID
	nextID++
	certificates = append(certificates, newCertificate)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newCertificate)
}

func GetAllCertificates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(certificates)
}

func UpdateCertificate(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var updatedCertificate Certificate
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&updatedCertificate); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i, certificate := range certificates {
		if certificate.ID == id {
			// Update the certificate data
			certificates[i] = updatedCertificate
			certificates[i].ID = id // Ensure the ID remains the same
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(certificates[i])
			return
		}
	}

	http.Error(w, "Certificate not found", http.StatusNotFound)
}

func main() {
	certificates = append(certificates, Certificate{ID: nextID, Name: "GoLang Certification", Owner: "John Doe"})
	nextID++

	r := mux.NewRouter()

	r.HandleFunc("/certificates/{id}", GetCertificateByID).Methods("GET") // Get certificate by ID
	r.HandleFunc("/certificates", CreateCertificate).Methods("POST")      // Create new certificate
	r.HandleFunc("/certificates", GetAllCertificates).Methods("GET")      // Get all certificates
	r.HandleFunc("/certificates/{id}", UpdateCertificate).Methods("PUT")  // Update certificate

	fmt.Println("Server is running on port 8000...")
	log.Fatal(http.ListenAndServe(":8000", r))
}
