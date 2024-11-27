package controller

import (
	"github.com/go-chi/chi/v5"
	"github.com/vlegro/backend/api/handlers"
	"github.com/vlegro/backend/api/service"
)

type CustomerController struct {
	customerHandler *handlers.CustomerHandler
}

func NewCustomerController(customerHandler *service.CustomerService) *CustomerController {
	return &CustomerController{customerHandler: handlers.NewCustomerHandler(customerHandler)}
}

func (cc *CustomerController) RestController() chi.Router {
	router := chi.NewRouter()

	// Add routes
	router.Get("/customers", cc.customerHandler.HandleGetByPrefix)
	router.Delete("/customers", cc.customerHandler.HandleDeleteByPrefix)

	return router
}
