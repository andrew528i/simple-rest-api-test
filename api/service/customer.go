package service

import (
	"fmt"
	"strings"

	"github.com/vlegro/backend/api/repository"
)

type CustomerService struct {
	customerRepository repository.CustomerRepository
}

func NewCustomerService(customerRepository repository.CustomerRepository) *CustomerService {
	return &CustomerService{customerRepository: customerRepository}
}

func (cs *CustomerService) Get(prefix []string) ([]repository.CustomerInfo, error) {
	// Validate input
	if len(prefix) == 0 {
		return nil, fmt.Errorf("prefix cannot be empty")
	}

	// Get customers by prefix
	customers, err := cs.customerRepository.GetByPrefix(prefix)
	if err != nil {
		return nil, fmt.Errorf("failed to get customers: %w", err)
	}

	return customers, nil
}

func (cs *CustomerService) Delete(prefix string) (repository.DeleteInfo, error) {
	// Validate input
	if prefix == "" {
		return repository.DeleteInfo{}, fmt.Errorf("prefix cannot be empty")
	}

	// Split prefix by comma and trim spaces
	prefixes := strings.Split(prefix, ",")
	for i := range prefixes {
		prefixes[i] = strings.TrimSpace(prefixes[i])
		if prefixes[i] == "" {
			return repository.DeleteInfo{}, fmt.Errorf("invalid empty prefix in the list")
		}
	}

	// Delete customers by prefix
	deleteInfo, err := cs.customerRepository.DeleteByPrefix(prefixes)
	if err != nil {
		return repository.DeleteInfo{}, fmt.Errorf("failed to delete customers: %w", err)
	}

	return deleteInfo, nil
}