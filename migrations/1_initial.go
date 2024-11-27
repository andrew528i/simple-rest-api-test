package main

import (
	"gorm.io/gorm"
)

// nolint:funlen
func init() {
	// nolint:misspell
	addMigration("1_initial.go",
		func(tx *gorm.DB) error {
			result := tx.Exec(
				`CREATE TABLE customer(
					id integer CONSTRAINT pk PRIMARY KEY,
					first_name varchar(100),
					last_name varchar(100),
					patronymic_name varchar(100),
					phone varchar(20),
					email varchar(200)
					);
			insert into customer (id, first_name, last_name, patronymic_name, phone, email) values (1, 'Клиент1', 'Клиентов1', 'Клиентович1', '77777777777', 'test1@test.ru');
			insert into customer (id, first_name, last_name, patronymic_name, phone, email) values (2, 'Клиент2', 'Клиентов2', 'Клиентович2', '77777777777', 'test2@test.ru');
			insert into customer (id, first_name, last_name, patronymic_name, phone, email) values (3, 'Клиент3', 'Клиентов3', 'Клиентович3', '77777777777', 'test3@test.ru');
			insert into customer (id, first_name, last_name, patronymic_name, phone, email) values (4, 'Клиент4', 'Клиентов4', 'Клиентович4', '77777777777', 'test4@test.ru');
			insert into customer (id, first_name, last_name, patronymic_name, phone, email) values (5, 'ДругойКлиент5', 'Клиентов5', 'Клиентович5', '77777777777', 'test5@test.ru');
		    `)
			return result.Error
		},
		func(tx *gorm.DB) error {
			result := tx.Exec(`
			-- You need to delete all the tables to cancel the migration.
            -- So keep it empty to prevent data loss.
		`)
			return result.Error
		},
	)
}
