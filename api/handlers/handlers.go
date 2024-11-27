package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/vlegro/backend/api/service"
)

type CustomerHandler struct {
	customerService *service.CustomerService
}

func NewCustomerHandler(customerService *service.CustomerService) *CustomerHandler {
	return &CustomerHandler{customerService: customerService}
}

func (ch *CustomerHandler) HandleDeleteByPrefix(w http.ResponseWriter, r *http.Request) {
	// Only allow DELETE method
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get prefix from query parameter
	prefix := r.URL.Query().Get("prefix")
	if prefix == "" {
		http.Error(w, "prefix parameter is required", http.StatusBadRequest)
		return
	}

	// Delete customers
	deleteInfo, err := ch.customerService.Delete(prefix)
	if err != nil {
		log.Printf("Error deleting customers: %v", err)
		http.Error(w, "Failed to delete customers", http.StatusInternalServerError)
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Write response
	if err := json.NewEncoder(w).Encode(deleteInfo); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

func (ch *CustomerHandler) HandleGetByPrefix(w http.ResponseWriter, r *http.Request) {
	// Only allow GET method
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get prefix from query parameter
	prefix := r.URL.Query().Get("prefix")
	if prefix == "" {
		http.Error(w, "prefix parameter is required", http.StatusBadRequest)
		return
	}

	// Split prefix by comma and trim spaces
	prefixes := strings.Split(prefix, ",")
	for i := range prefixes {
		prefixes[i] = strings.TrimSpace(prefixes[i])
	}

	// Get customers
	customers, err := ch.customerService.Get(prefixes)
	if err != nil {
		log.Printf("Error getting customers: %v", err)
		http.Error(w, "Failed to get customers", http.StatusInternalServerError)
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Write response
	if err := json.NewEncoder(w).Encode(customers); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}