package models

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/naman-dave/gqlgo/internal/db"
)

type CarUsage struct {
	ID          int64  `json:"id"`
	CarUniqueID string `json:"car_unique_id"`
	UserID      string `json:"user_id"`
	BookedTill  string `json:"booked_till"`
	RetunedDate string `json:"retuned_date"`
}

func (cu *CarUsage) BookCar() (int64, error) {

	car, err := GetCarByCarIdentifier(cu.CarUniqueID)
	if err != nil {
		return 0, fmt.Errorf("car %s not found", cu.CarUniqueID)
	}

	if car.TotalCar == car.TotalInUse {
		return 0, fmt.Errorf("all cars %s is already in use", car.CarIdentifier)
	}

	insertSQL := ` INSERT INTO CarUsages (carunidueid, userid, bookedtill) Values ($1, $2, $3) RETURNING id`

	dbTX, err := db.Db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return 0, err
	}

	err = dbTX.QueryRow(insertSQL, cu.CarUniqueID, cu.UserID, cu.BookedTill).Scan(&cu.ID)
	if err != nil {
		log.Println("carusgae.BookCar(): ", err)
		return 0, err
	}

	err = car.bookCar()
	if err != nil {
		dbTX.Rollback()
		log.Println("carusgae.BookCar(): ", err)
		return 0, err
	}

	dbTX.Commit()
	return cu.ID, nil
}

func (cu *CarUsage) ReturnCar(billno int64) error {
	car, err := GetCarByCarIdentifier(cu.CarUniqueID)
	if err != nil {
		return fmt.Errorf("car %s not found", cu.CarUniqueID)
	}

	updateSQL := ` Update CarUsages
	    Set returndate = $1
		WHERE id = $2 `

	dbTX, err := db.Db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		log.Println("carusgaes.ReturnCar(): ", err)
		return err
	}

	_, err = dbTX.Exec(updateSQL, cu.BookedTill, billno)
	if err != nil {
		log.Println("carusgaes.ReturnCar(): ", err)
		return err
	}

	err = car.returnCar()
	if err != nil {
		dbTX.Rollback()
		log.Println("carusgaes.ReturnCar(): ", err)
		return err
	}

	dbTX.Commit()
	return nil
}

func GetCarUsage(billno int64) (*CarUsage, error) {
	cu := CarUsage{}

	returnDate := sql.NullString{}

	selectSQL := ` SELECT id, carunidueid, userid, bookedtill, DATE(returndate) FROM carusages where id = $1 `

	err := db.Db.QueryRow(selectSQL, billno).Scan(&cu.ID, &cu.CarUniqueID, &cu.UserID, &cu.BookedTill, &returnDate)
	if err != nil {
		log.Println("GetCarUsage():", err)
		return nil, err
	}

	cu.RetunedDate = returnDate.String
	return &cu, nil
}
