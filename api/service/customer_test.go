package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vlegro/backend/api/repository"
)

// MockCustomerRepository is a mock implementation of CustomerRepository
type MockCustomerRepository struct {
	mock.Mock
}

func (m *MockCustomerRepository) GetByPrefix(prefix []string) ([]repository.CustomerInfo, error) {
	args := m.Called(prefix)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]repository.CustomerInfo), args.Error(1)
}

func (m *MockCustomerRepository) DeleteByPrefix(prefix []string) (repository.DeleteInfo, error) {
	args := m.Called(prefix)
	return args.Get(0).(repository.DeleteInfo), args.Error(1)
}

func TestCustomerService_Get(t *testing.T) {
	mockRepo := new(MockCustomerRepository)
	service := NewCustomerService(mockRepo)

	tests := []struct {
		name           string
		prefix         []string
		mockReturn     []repository.CustomerInfo
		mockError      error
		expectedResult []repository.CustomerInfo
		expectedError  bool
	}{
		{
			name:   "successful get",
			prefix: []string{"Клиент"},
			mockReturn: []repository.CustomerInfo{
				{
					Id:        1,
					FirstName: strPtr("Клиент1"),
					LastName:  strPtr("Клиентов1"),
				},
			},
			mockError:      nil,
			expectedResult: []repository.CustomerInfo{{Id: 1, FirstName: strPtr("Клиент1"), LastName: strPtr("Клиентов1")}},
			expectedError:  false,
		},
		{
			name:           "empty prefix",
			prefix:         []string{},
			mockReturn:     nil,
			mockError:      nil,
			expectedResult: nil,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.prefix) > 0 {
				mockRepo.On("GetByPrefix", tt.prefix).Return(tt.mockReturn, tt.mockError)
			}

			result, err := service.Get(tt.prefix)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestCustomerService_Delete(t *testing.T) {
	mockRepo := new(MockCustomerRepository)
	service := NewCustomerService(mockRepo)

	tests := []struct {
		name           string
		prefix         string
		mockReturn     repository.DeleteInfo
		mockError      error
		expectedResult repository.DeleteInfo
		expectedError  bool
	}{
		{
			name:           "successful delete",
			prefix:         "Клиент",
			mockReturn:     repository.DeleteInfo{Count: 1, Ids: []int{1}},
			mockError:      nil,
			expectedResult: repository.DeleteInfo{Count: 1, Ids: []int{1}},
			expectedError:  false,
		},
		{
			name:           "empty prefix",
			prefix:         "",
			mockReturn:     repository.DeleteInfo{},
			mockError:      nil,
			expectedResult: repository.DeleteInfo{},
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.prefix != "" {
				mockRepo.On("DeleteByPrefix", []string{tt.prefix}).Return(tt.mockReturn, tt.mockError)
			}

			result, err := service.Delete(tt.prefix)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func strPtr(s string) *string {
	return &s
}
