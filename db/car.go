package db

import (
	"catalog/app"
	"catalog/models"
	"context"
	"database/sql"
	"errors"
	"github.com/gofrs/uuid"
	jsonit "github.com/json-iterator/go"
	"github.com/lib/pq"
)

// External errors

var NotFound = errors.New("not found")
var AlreadyExists = errors.New("already exists")

// Internal errors
var uniqueViolation = pq.ErrorCode("23505")

func isUniqueViolation(err error) bool {
	var pgErr *pq.Error
	ok := errors.As(err, &pgErr)
	return ok && pgErr.Code == uniqueViolation
}

// db functions

func GetCars(app *app.App, filter models.Car, offset int, limit int) ([]models.Car, error) {
	tx, err := app.DB().BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	rows, err := tx.Query(
		`SELECT * FROM cars
			OFFSET $1
			LIMIT $2`,
		offset,
		limit,
	)
	if err != nil {
		return nil, err
	}

	var cars []models.Car

	for rows.Next() {
		var car models.Car
		var owner []byte
		err = rows.Scan(&car.ID, &car.RegNum, &car.Mark, &car.Model, &car.Year, &owner)
		if err != nil {
			return nil, err
		}
		err = jsonit.Unmarshal(owner, &car.Owner)
		if err != nil {
			return nil, err
		}
		cars = append(cars, car)
	}

	var filteredCars []models.Car
	for _, car := range cars {
		if filter.ID != uuid.Nil && car.ID != filter.ID {
			continue
		}
		if filter.RegNum != "" && car.RegNum != filter.RegNum {
			continue
		}
		if filter.Mark != "" && car.Mark != filter.Mark {
			continue
		}
		if filter.Model != "" && car.Model != filter.Model {
			continue
		}
		if filter.Year != 0 && car.Year != filter.Year {
			continue
		}

		filteredCars = append(filteredCars, car)
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return filteredCars, nil
}

func AddCars(app *app.App, cars []*models.Car) error {
	tx, err := app.DB().BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, car := range cars {
		car.ID = uuid.Must(uuid.NewV4())
		ownerBytes, err := jsonit.Marshal(car.Owner)
		if err != nil {
			return err
		}
		_, err = tx.Exec(
			`INSERT INTO cars (id, regNum, mark, model, year, owner)
				VALUES ($1, $2, $3, $4, $5, $6)`,
			car.ID, car.RegNum, car.Mark, car.Model, car.Year, ownerBytes,
		)

		if isUniqueViolation(err) {
			return AlreadyExists
		} else if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func UpdateCars(app *app.App, cars []*models.Car) error {
	tx, err := app.DB().BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, car := range cars {
		ownerBytes, err := jsonit.Marshal(car.Owner)
		if err != nil {
			return err
		}
		row := tx.QueryRow(
			`UPDATE cars
				SET regNum = $2, mark = $3, model = $4, year = $5, owner = $6
				WHERE id = $1
				RETURNING id`,
			car.ID, car.RegNum, car.Mark, car.Model, car.Year, ownerBytes,
		)
		var discardID uuid.UUID
		err = row.Scan(&discardID)
		if isUniqueViolation(err) {
			return AlreadyExists
		} else if errors.Is(err, sql.ErrNoRows) {
			return NotFound
		} else if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func DeleteCarById(app *app.App, id uuid.UUID) error {
	tx, err := app.DB().BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	row := tx.QueryRow(
		`DELETE FROM cars
			WHERE id = $1
			RETURNING id`,
		id,
	)
	var discardID uuid.UUID
	err = row.Scan(&discardID)
	if errors.Is(err, sql.ErrNoRows) {
		return NotFound
	} else if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
