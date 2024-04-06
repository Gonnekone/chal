package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Gonnekone/challenge/internal/domain/models"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
}

func New(log *slog.Logger, dbURL string) (*Storage, error) {
	const op = "storage.postgres.New"

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("creating table 'people'")

	_, err = pool.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS people (
			id SERIAL PRIMARY KEY,
			name TEXT UNIQUE NOT NULL,
			surname TEXT UNIQUE NOT NULL,
			patronymic TEXT UNIQUE
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("creating table 'car'")

	_, err = pool.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS car (
			id SERIAL PRIMARY KEY,
			regNum TEXT UNIQUE NOT NULL,
			mark TEXT NOT NULL,
			model TEXT NOT NULL,
			year INT,
			owner_id INT REFERENCES people(id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: pool}, nil
}

func (s *Storage) GetCar(ctx context.Context, log *slog.Logger, filters map[string]interface{}, limit, offset int) ([]models.Car, error) {
    const op = "postgres.getCar"

	log = log.With(
		slog.String("op", op),
	)
	
	query := "SELECT c.regNum, c.mark, c.model, c.year, p.name, p.surname, p.patronymic " +
        "FROM car c JOIN people p ON c.owner_id = p.id"

	if len(filters) != 0 {
		query += " WHERE"
		for k, v := range filters {
			if k != "year" {
				query += fmt.Sprintf(" %s = '%s' AND", k, v)
			} else {
				query += fmt.Sprintf(" %s = %d AND", k, v)
			}
		}

		query = query[:len(query)-3]
	}

    query += " LIMIT $1::bigint OFFSET $2::bigint"
	
	log.Debug("getting cars", slog.String("query", query))
	
    rows, err := s.db.Query(ctx, query, limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var cars []models.Car
    for rows.Next() {
        var car models.Car
        var owner models.People
        err := rows.Scan(&car.RegNum, &car.Mark, &car.Model, &car.Year, &owner.Name, &owner.Surname, &owner.Patronymic)
        if err != nil {
            return nil, err
        }
        car.Owner = owner
        cars = append(cars, car)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return cars, nil
}


func (s *Storage) SaveOwner(ctx context.Context, owner models.People) (int64, error) {
	var id int64
	err := s.db.QueryRow(ctx,
		"INSERT INTO people(name, surname, patronymic) VALUES($1, $2, $3) RETURNING id",
		owner.Name, owner.Surname, owner.Patronymic).Scan(&id)

	if err != nil {
		pgErr, ok := err.(*pgconn.PgError)
		if ok && pgErr.Code == "23505" {
			err := s.db.QueryRow(ctx,
				"SELECT id FROM people WHERE name = $1 AND surname = $2",
				owner.Name, owner.Surname).Scan(&id)
			if err != nil {
				return 0, err
			}
			return id, nil
		}
		return 0, err
	}
	return id, nil
}


func (s *Storage) DeleteCar(ctx context.Context, log *slog.Logger, id int64) error {
	const op = "postgres.deleteCar"

	log = log.With(
		slog.String("op", op),
	)

	log.Debug(fmt.Sprintf("deleting car with id = %d", id), slog.String("query", "DELETE FROM car WHERE id = $1"))

	_, err := s.db.Exec(ctx, "DELETE FROM car WHERE id = $1", id)
	return err
}

func (s *Storage) UpdateCar(ctx context.Context, log *slog.Logger, id int64, updates map[string]interface{}) error {
	const op = "postgres.updateCar"

	log = log.With(
		slog.String("op", op),
	)

	query := "UPDATE car SET"

	if len(updates) != 0 {
		for k, v := range updates {
			if k != "year" {
				query += fmt.Sprintf(" %s = '%s',", k, v)
			} else {
				query += fmt.Sprintf(" %s = %d,", k, v)
			}
		}
	}

	query = query[:len(query)-1]
	query += " WHERE id = $1"
	
	log.Debug(fmt.Sprintf("updating car with id = %d", id), slog.String("query", query))

	_, err := s.db.Exec(ctx, query, id)
	return err
}

func (s *Storage) SaveCar(ctx context.Context, log *slog.Logger, regNums []string) error {
	const op = "postgres.saveCar"

	log = log.With(
		slog.String("op", op),
	)

	query := "INSERT INTO car (regNum, mark, model, year, owner_id) VALUES ($1, $2, $3, $4, $5)"

	for _, regNum := range regNums {
		car, err := fetchCarInfo(ctx, log, regNum)
		if err != nil {
			return fmt.Errorf("error fetching car info for regNum %s: %w", regNum, err)
		}

		ownerId, err := s.SaveOwner(ctx, car.Owner)
		if err != nil {
			return fmt.Errorf("error saving owner for regNum %s: %w", regNum, err)
		}

		log.Debug(fmt.Sprintf("saving car with regNum = %s", regNum), slog.Any("car", car), slog.String("query", query))


		_, err = s.db.Exec(ctx, query, car.RegNum, car.Mark, car.Model, car.Year, ownerId)
		if err != nil {
			return fmt.Errorf("error inserting car info for regNum %s: %w", regNum, err)
		}
	}

	return nil
}

func fetchCarInfo(ctx context.Context, log *slog.Logger, regNum string) (models.Car, error) {
	const op = "postgres.fetchCarInfo"

	log = log.With(
		slog.String("op", op),
	)

	apiURL := fmt.Sprintf("http://localhost:8086/car?regNum=%s", regNum)

	log.Debug(fmt.Sprintf("fetching car with regNum = %s", regNum), slog.String("apiURL", apiURL))

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return models.Car{}, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return models.Car{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.Car{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var car models.Car
	if err := json.NewDecoder(resp.Body).Decode(&car); err != nil {
		return models.Car{}, err
	}

	log.Debug("car fetched", slog.Any("car", car))

	return car, nil
}

// func fetchCarInfoTest(ctx context.Context, regNum string) (models.Car, error) {
// 	cache := make(map[string]models.Car)

// 	car1 := models.Car {
// 		RegNum: "X123XX150",
// 		Mark:   "Lada",
// 		Model:  "Vesta",
// 		Owner: models.People {
// 			Name:      "John",
// 			Surname:   "Doe",
// 			Patronymic: "Smith",
// 		},
// 	}
	
// 	car2 := models.Car {
// 		RegNum: "A456BC789",
// 		Mark:   "Toyota",
// 		Model:  "Corolla",
// 		Year:   2015,
// 		Owner: models.People {
// 			Name:    "Alice",
// 			Surname: "Johnson",
// 		},
// 	}
	
// 	car3 := models.Car {
// 		RegNum: "H789GF123",
// 		Mark:   "BMW",
// 		Model:  "X5",
// 		Year:   2019,
// 		Owner: models.People {
// 			Name:    "Bob",
// 			Surname: "Brown",
// 			Patronymic: "Lee",
// 		},
// 	}

// 	car4 := models.Car {
// 		RegNum: "Z456BC789",
// 		Mark:   "Toyota",
// 		Model:  "Mark 2",
// 		Year:   1999,
// 		Owner: models.People{
// 			Name:    "Alice",
// 			Surname: "Johnson",
// 		},
// 	}
	
	
// 	cache[car1.RegNum] = car1
// 	cache[car2.RegNum] = car2
// 	cache[car3.RegNum] = car3
// 	cache[car4.RegNum] = car4

// 	return cache[regNum], nil
// }
