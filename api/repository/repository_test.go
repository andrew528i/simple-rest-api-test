package repository

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sql.DB {
	// Connect to the database
	db, err := sql.Open("postgres", "postgres://backend:backend@localhost:5555/backend?sslmode=disable")
	require.NoError(t, err)
	require.NoError(t, db.Ping())
	return db
}

func TestCustomerRepository_GetByPrefix(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewCustomerRepositoryImpl(db)

	tests := []struct {
		name           string
		prefixes       []string
		expectedCount  int
		expectedNames  []string
		expectedError  bool
	}{
		{
			name:           "single prefix match",
			prefixes:       []string{"Клиент"},
			expectedCount:  4,
			expectedNames:  []string{"Клиент1", "Клиент2", "Клиент3", "Клиент4"},
			expectedError:  false,
		},
		{
			name:           "multiple prefixes match",
			prefixes:       []string{"Клиент", "Другой"},
			expectedCount:  5,
			expectedNames:  []string{"Клиент1", "Клиент2", "Клиент3", "Клиент4", "ДругойКлиент5"},
			expectedError:  false,
		},
		{
			name:           "no matches",
			prefixes:       []string{"NonExistent"},
			expectedCount:  0,
			expectedNames:  []string{},
			expectedError:  false,
		},
		{
			name:           "empty prefix list",
			prefixes:       []string{},
			expectedCount:  0,
			expectedNames:  []string{},
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			customers, err := repo.GetByPrefix(tt.prefixes)

			if tt.expectedError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Len(t, customers, tt.expectedCount)

			if tt.expectedCount > 0 {
				names := make([]string, len(customers))
				for i, c := range customers {
					names[i] = *c.FirstName
				}
				assert.ElementsMatch(t, tt.expectedNames, names)
			}
		})
	}
}

func TestCustomerRepository_DeleteByPrefix(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewCustomerRepositoryImpl(db)

	tests := []struct {
		name           string
		prefixes       []string
		expectedCount  int
		expectedError  bool
		setup         func(t *testing.T, db *sql.DB)
		cleanup       func(t *testing.T, db *sql.DB)
	}{
		{
			name:          "delete single prefix",
			prefixes:      []string{"Клиент"},
			expectedCount: 4,
			expectedError: false,
			setup: func(t *testing.T, db *sql.DB) {
				// No setup needed, using initial migration data
			},
			cleanup: func(t *testing.T, db *sql.DB) {
				// Restore original data
				_, err := db.Exec(`
					INSERT INTO customer (id, first_name, last_name, patronymic_name, phone, email)
					VALUES 
						(1, 'Клиент1', 'Клиентов1', 'Клиентович1', '77777777777', 'test1@test.ru'),
						(2, 'Клиент2', 'Клиентов2', 'Клиентович2', '77777777777', 'test2@test.ru'),
						(3, 'Клиент3', 'Клиентов3', 'Клиентович3', '77777777777', 'test3@test.ru'),
						(4, 'Клиент4', 'Клиентов4', 'Клиентович4', '77777777777', 'test4@test.ru')
				`)
				require.NoError(t, err)
			},
		},
		{
			name:          "delete multiple prefixes",
			prefixes:      []string{"Клиент", "Другой"},
			expectedCount: 5,
			expectedError: false,
			setup: func(t *testing.T, db *sql.DB) {
				// No setup needed, using initial migration data
			},
			cleanup: func(t *testing.T, db *sql.DB) {
				// Restore all original data
				_, err := db.Exec(`
					INSERT INTO customer (id, first_name, last_name, patronymic_name, phone, email)
					VALUES 
						(1, 'Клиент1', 'Клиентов1', 'Клиентович1', '77777777777', 'test1@test.ru'),
						(2, 'Клиент2', 'Клиентов2', 'Клиентович2', '77777777777', 'test2@test.ru'),
						(3, 'Клиент3', 'Клиентов3', 'Клиентович3', '77777777777', 'test3@test.ru'),
						(4, 'Клиент4', 'Клиентов4', 'Клиентович4', '77777777777', 'test4@test.ru'),
						(5, 'ДругойКлиент5', 'Клиентов5', 'Клиентович5', '77777777777', 'test5@test.ru')
				`)
				require.NoError(t, err)
			},
		},
		{
			name:          "delete non-existent prefix",
			prefixes:      []string{"NonExistent"},
			expectedCount: 0,
			expectedError: false,
			setup: func(t *testing.T, db *sql.DB) {
				// No setup needed
			},
			cleanup: func(t *testing.T, db *sql.DB) {
				// No cleanup needed
			},
		},
		{
			name:          "empty prefix list",
			prefixes:      []string{},
			expectedCount: 0,
			expectedError: true,
			setup: func(t *testing.T, db *sql.DB) {
				// No setup needed
			},
			cleanup: func(t *testing.T, db *sql.DB) {
				// No cleanup needed
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Run setup if provided
			if tt.setup != nil {
				tt.setup(t, db)
			}

			// Run test
			result, err := repo.DeleteByPrefix(tt.prefixes)

			// Verify results
			if tt.expectedError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expectedCount, result.Count)
			assert.Len(t, result.Ids, tt.expectedCount)

			// Run cleanup if provided
			if tt.cleanup != nil {
				tt.cleanup(t, db)
			}
		})
	}
}
