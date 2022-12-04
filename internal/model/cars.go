package models

import (
	"fmt"
	"log"

	"github.com/naman-dave/gqlgo/graph/model"
	"github.com/naman-dave/gqlgo/internal/db"
)

type Car struct {
	ID                int    `json:"id"`
	CarIdentifier     string `json:"caridentifier"`
	Name              string `json:"name"`
	DateOfManufacture string `json:"date_of_manufacture"`
	TotalCar          int    `json:"total_car"`
	TotalInUse        int    `json:"total_in_use"`
}

func (c *Car) Insert() (string, error) {

	inserSQL := ` INSERT INTO Cars (
		CarIdentifier,
		Name,
		DateOfManufacture,
        Total,
        TotalInUse
	) Values (
		$1,
		$2,
        $3,
        $4,
		$5,
        $6,
        $7
	)`

	_, err := db.Db.Exec(inserSQL, c.CarIdentifier, c.Name, c.DateOfManufacture, c.TotalCar, c.TotalInUse)
	if err != nil {
		log.Println("car.Insert():", err)
		return "", err
	}

	return c.CarIdentifier, nil
}

func (c *Car) BookCar() error {
	updateSQL := `
	    UPDATE Cars
        SET
			TotalInUse = $1 
		WHERE 
			CarIdentifier = $2 `

	_, err := db.Db.Exec(updateSQL, c.TotalInUse+1, c.CarIdentifier)
	if err != nil {
		log.Println("car.BookCar():", err)
		return err
	}

	return nil
}

func (c *Car) ReturnCar() error {
	updateSQL := `
	    UPDATE Cars
        SET
			TotalInUse = $1 
		WHERE 
			CarIdentifier = $2 `

	if c.TotalInUse == 0 {
		return fmt.Errorf("there are no cars in use")
	}

	_, err := db.Db.Exec(updateSQL, c.TotalInUse-1, c.CarIdentifier)
	if err != nil {
		log.Println("car.BookCar():", err)
		return err
	}

	return nil
}

func GetCarByCarIdentifier(id string) (*Car, error) {
	car := Car{}

	selectSQL := ` SELECT 
	    CarIdentifier,
		Name,
		DATE(DateOfManufacture),
        Total,
        TotalInUse
		FROM Cars 
		WHERE CarIdentifier = $1`

	err := db.Db.QueryRow(selectSQL, id).Scan(&car.CarIdentifier, &car.Name, &car.DateOfManufacture, &car.TotalCar, &car.TotalInUse)
	if err != nil {
		log.Println("GetCarByCarIdentifier():", err)
		return nil, err
	}

	return &car, nil
}

func GetCars(cars *[]*model.Car) error {

	selectSQL := ` SELECT 
	    CarIdentifier,
		Name,
		DATE(DateOfManufacture),
        Total,
        TotalInUse
		FROM Cars `

	rows, err := db.Db.Query(selectSQL)
	if err != nil {
		log.Println(err)
		return err
	}
	defer rows.Close()

	tempcar := *cars

	for rows.Next() {
		car := &model.Car{}
		err := rows.Scan(&car.Caridentifier, &car.Modal, &car.Dateofmanufacture, &car.Totalcar, &car.Totalinuse)
		if err != nil {
			log.Println(err)
			return err
		}
		tempcar = append(tempcar, car)
	}

	*cars = tempcar
	return nil
}
