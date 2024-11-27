package repository

import (
	"database/sql"
	"fmt"
	"strings"
)

type CustomerRepositoryImpl struct {
	dbConnection *sql.DB
}

func NewCustomerRepositoryImpl(dbConnection *sql.DB) *CustomerRepositoryImpl {
	return &CustomerRepositoryImpl{
		dbConnection: dbConnection,
	}
}

func (c *CustomerRepositoryImpl) GetByPrefix(prefixes []string) ([]CustomerInfo, error) {
	// Build the WHERE clause for multiple prefixes
	conditions := make([]string, len(prefixes))
	args := make([]interface{}, len(prefixes))
	for i, prefix := range prefixes {
		conditions[i] = fmt.Sprintf("first_name LIKE $%d", i+1)
		args[i] = prefix + "%"
	}
	whereClause := strings.Join(conditions, " OR ")

	// Prepare the query
	query := fmt.Sprintf(`
		SELECT id, first_name, last_name, patronymic_name, phone, email 
		FROM customer 
		WHERE %s
		ORDER BY id`, whereClause)

	// Execute the query
	rows, err := c.dbConnection.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	// Process results
	var customers []CustomerInfo
	for rows.Next() {
		var customer CustomerInfo
		var firstName, lastName, patronymicName, phone, email sql.NullString
		
		err := rows.Scan(
			&customer.Id,
			&firstName,
			&lastName,
			&patronymicName,
			&phone,
			&email,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Handle nullable fields
		if firstName.Valid {
			customer.FirstName = &firstName.String
		}
		if lastName.Valid {
			customer.LastName = &lastName.String
		}
		if patronymicName.Valid {
			customer.PatronymicName = &patronymicName.String
		}
		if phone.Valid {
			customer.Phone = &phone.String
		}
		if email.Valid {
			customer.Email = &email.String
		}

		customers = append(customers, customer)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return customers, nil
}

func (c *CustomerRepositoryImpl) DeleteByPrefix(prefixes []string) (DeleteInfo, error) {
	// Start a transaction
	tx, err := c.dbConnection.Begin()
	if err != nil {
		return DeleteInfo{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Will be ignored if transaction is committed

	// Build the WHERE clause for multiple prefixes
	conditions := make([]string, len(prefixes))
	args := make([]interface{}, len(prefixes))
	for i, prefix := range prefixes {
		conditions[i] = fmt.Sprintf("first_name LIKE $%d", i+1)
		args[i] = prefix + "%"
	}
	whereClause := strings.Join(conditions, " OR ")

	// First get the IDs of customers to be deleted
	selectQuery := fmt.Sprintf("SELECT id FROM customer WHERE %s", whereClause)
	rows, err := tx.Query(selectQuery, args...)
	if err != nil {
		return DeleteInfo{}, fmt.Errorf("failed to query customers: %w", err)
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return DeleteInfo{}, fmt.Errorf("failed to scan id: %w", err)
		}
		ids = append(ids, id)
	}
	if err = rows.Err(); err != nil {
		return DeleteInfo{}, fmt.Errorf("error iterating rows: %w", err)
	}

	// If no customers found, return early
	if len(ids) == 0 {
		return DeleteInfo{Count: 0, Ids: []int{}}, nil
	}

	// Delete the customers
	deleteQuery := fmt.Sprintf("DELETE FROM customer WHERE %s", whereClause)
	result, err := tx.Exec(deleteQuery, args...)
	if err != nil {
		return DeleteInfo{}, fmt.Errorf("failed to delete customers: %w", err)
	}

	// Get number of affected rows
	count, err := result.RowsAffected()
	if err != nil {
		return DeleteInfo{}, fmt.Errorf("failed to get affected rows: %w", err)
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return DeleteInfo{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return DeleteInfo{
		Count: int(count),
		Ids:   ids,
	}, nil
}
