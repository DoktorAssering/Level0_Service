package handler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"mainservice/ServiceApp/service"
	"net/http"
	"strconv"
)

type DatabaseHandler struct {
	databaseService *service.DatabaseService
}

func NewDatabaseHandler(databaseService *service.DatabaseService) *DatabaseHandler {
	return &DatabaseHandler{
		databaseService: databaseService,
	}
}

func (h *DatabaseHandler) Index(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("../../web/index.html")
	if err != nil {
		http.Error(w, "Template parsing error", http.StatusInternalServerError)
		log.Println("Template parsing error:", err)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Template parsing error", http.StatusInternalServerError)
		log.Println("Template parsing error:", err)
		return
	}
}

func (h *DatabaseHandler) GetInfo(w http.ResponseWriter, r *http.Request) {
	numberStr := r.FormValue("number")
	number, err := strconv.Atoi(numberStr)
	if err != nil {
		http.Error(w, "Invalid number", http.StatusBadRequest)
		return
	}

	jsonData, err := h.databaseService.GetInfo(number)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tmpl, err := template.ParseFiles("../../web/index.html")
	if err != nil {
		http.Error(w, "Template parsing error", http.StatusInternalServerError)
		log.Println("Template parsing error:", err)
		return
	}

	data := struct {
		Info   string
		Number int
	}{
		Info:   jsonData,
		Number: number,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Template display error", http.StatusInternalServerError)
		log.Println("Template display error:", err)
		return
	}
}

func (h *DatabaseHandler) AddData(w http.ResponseWriter, r *http.Request) {
	jsonData := r.FormValue("jsonData")

	id, err := h.databaseService.AddData(jsonData)
	if err != nil {
		http.Error(w, "Error adding data", http.StatusInternalServerError)
		log.Printf("Error adding data: %v", err)
		return
	}

	fmt.Fprintf(w, "New data added, ID: %d", id)
}

func (h *DatabaseHandler) GetAllIDsHandler(w http.ResponseWriter, r *http.Request) {
	ids, err := h.databaseService.GetAllIDs()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ids)
}
