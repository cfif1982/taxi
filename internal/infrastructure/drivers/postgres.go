package drivers

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/cfif1982/taxi/internal/domain/drivers"
	"github.com/cfif1982/taxi/pkg/logger"
	"github.com/google/uuid"
	"github.com/pressly/goose/v3"

	"github.com/jackc/pgerrcode"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/jackc/pgx/v5/pgconn"
)

// postgres репозиторий
type PostgresRepository struct {
	db *sql.DB
}

// Создаем репозиторий БД
func NewPostgresRepository(ctx context.Context, databaseDSN string, logger *logger.Logger) (*PostgresRepository, error) {

	db, err := sql.Open("pgx", databaseDSN)
	if err != nil {
		return nil, err
	}

	// QUESTION: нужно ли здесь пинговать БД для проверки ее доступности? Это нормальная практика?
	// создаю контекст для пинга
	// ctx2, cancel2 := context.WithTimeout(ctx, 1*time.Second)
	// defer cancel2()

	// пингую БД. Если не отвечает, то возвращаю ошибку
	// if err = db.PingContext(ctx2); err != nil {
	// 	return nil, err
	// }

	// начинаю миграцию
	logger.Info("Start migrating database")

	if err := goose.SetDialect("postgres"); err != nil {
		logger.Info(err.Error())
	}

	// узнаю текущую папку, чтобы передать путь к папке с миграциями
	ex, err := os.Executable()
	if err != nil {
		logger.Info(err.Error())
	}
	exPath := filepath.Dir(ex)

	exPath = exPath + "/migrations"

	err = goose.Up(db, exPath)
	if err != nil {
		logger.Info(err.Error() + ": " + exPath)
	}

	logger.Info("migrating database finished")

	return &PostgresRepository{
		db: db,
	}, nil
}

// добавить водителя
func (r *PostgresRepository) AddDriver(driver *drivers.Driver) error {

	// создаю контекст для запроса
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "INSERT INTO drivers(id, route_id, telephone, name, password, balance, last_paid_date) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	_, err := r.db.ExecContext(ctx, query,
		driver.ID(), driver.RouteID(), driver.Telephone(),
		driver.Name(), driver.Password(), driver.Balance(),
		driver.LastPaidDate())
	if err != nil {
		// проверяем ошибку на предмет вставки водителя с телефоном, которое уже есть в БД
		// создаем объект *pgconn.PgError - в нем будет храниться код ошибки из БД
		var pgErr *pgconn.PgError

		// преобразуем ошибку к типу pgconn.PgError
		if errors.As(err, &pgErr) {
			// если ошибка- запись существует, то возвращаем эту ошибку
			if pgErr.Code == pgerrcode.UniqueViolation {
				return drivers.ErrTelephoneAlreadyExist
			}
		} else {
			return err
		}
	}

	return nil
}

// Найти водителя по телефону
func (r *PostgresRepository) GetDriverByTelephone(telephone string) (*drivers.Driver, error) {

	// создаю контекст для запроса
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT id, route_id, name, password, balance, last_paid_date FROM drivers WHERE telephone=$1"
	row := r.db.QueryRowContext(ctx, query, telephone)

	// в эту переменную будет сканиться результат запроса
	var id, route_id uuid.UUID
	var name, password string
	var balance int
	var lastPaidDate time.Time

	err := row.Scan(&id, &route_id, &name, &password, &balance, &lastPaidDate)

	if err != nil {
		return nil, err
	}

	// создаем водителя и возвращаем его
	driver, err := drivers.NewDriver(id, route_id, telephone, name, password, balance, lastPaidDate)

	if err != nil {
		return nil, err
	}

	return driver, nil
}

// Найти водителя по id
func (r *PostgresRepository) GetDriverByID(id uuid.UUID) (*drivers.Driver, error) {

	// создаю контекст для запроса
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT route_id, name, password, telephone, balance, last_paid_date FROM drivers WHERE id=$1"
	row := r.db.QueryRowContext(ctx, query, id)

	// в эту переменную будет сканиться результат запроса
	var route_id uuid.UUID
	var name, password, telephone string
	var balance int
	var lastPaidDate time.Time

	err := row.Scan(&route_id, &name, &password, &telephone, &balance, &lastPaidDate)

	if err != nil {
		return nil, err
	}

	// создаем водителя и возвращаем его
	driver, err := drivers.NewDriver(id, route_id, telephone, name, password, balance, lastPaidDate)

	if err != nil {
		return nil, err
	}

	return driver, nil
}

// сохранить водителя
func (r *PostgresRepository) SaveDriver(driver *drivers.Driver) error {

	// создаю контекст для запроса
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "UPDATE drivers SET route_id=$1, name=$2, telephone=$3, password=$4, balance=$5, last_paid_date=$6 WHERE id=$7"
	_, err := r.db.ExecContext(ctx, query,
		driver.RouteID(), driver.Name(), driver.Telephone(),
		driver.Password(), driver.Balance(), driver.LastPaidDate(),
		driver.ID())

	if err != nil {
		return err
	}

	return nil
}
