package main

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	// реализуйте добавление строки в таблицу parcel, используйте данные из переменной p
	query := "INSERT INTO parcel (client, status, address, created_at) VALUES (?, ?, ?, ?)"
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(p.Client, p.Status, p.Address, p.CreatedAt)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	// верните идентификатор последней добавленной записи
	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка
	query := "SELECT * FROM parcel WHERE number = ?"
	row := s.db.QueryRow(query, number)
	// заполните объект Parcel данными из таблицы
	p := Parcel{}
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return Parcel{}, err
		}
		return Parcel{}, err
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк
	query := "SELECT * FROM parcel WHERE client = ?"
	rows, err := s.db.Query(query, client)
	if err != nil {
		return nil, fmt.Errorf("Error getting data about client with id %d: %v", client, err)
	}
	defer rows.Close()

	// заполните срез Parcel данными из таблицы
	var res []Parcel
	for rows.Next() {
		var p Parcel
		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	query := "UPDATE parcel SET status = ? WHERE number = ?"
	res, err := s.db.Exec(query, status, number)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("Error updating status of parcel: %v", err)
	}

	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	query := "UPDATE parcel SET address = ? WHERE number = ?"
	parcel, err := s.Get(number)
	if err != nil {
		return err
	}
	if parcel.Status != ParcelStatusRegistered {
		return fmt.Errorf("Cannot change address for parcel with status(1) %s", parcel.Status)
	}
	res, err := s.db.Exec(query, address, number)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("Incorrect status or parcel not found")
	}
	return nil
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	parcel, err := s.Get(number)
	if err != nil {
		return err
	}
	if parcel.Status != ParcelStatusRegistered {
		return fmt.Errorf("Cannot change address for parcel with status(2) %s", parcel.Status)
	}
	query := "DELETE FROM parcel WHERE number = ?"
	res, err := s.db.Exec(query, number)
	if err != nil {
		return fmt.Errorf("Failed to delete parcel: %s", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("Incorrect status or parcel not found")
	}
	return nil
}
