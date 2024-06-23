package routes

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/cfif1982/taxi/pkg/logger"
	"github.com/google/uuid"
	"github.com/pressly/goose/v3"

	"github.com/jackc/pgerrcode"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/cfif1982/taxi/internal/domain/routes"
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

	// начинаю миграцию
	// Т.к. делаю миграцию, то не нужно пинговать базу
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

// Получить все маршруты
func (r *PostgresRepository) GetAllRoutes() (*[]routes.Route, error) {

	// массив полученных маршрутов
	arrRoutes := make([]routes.Route, 0)

	// создаю контекст для запроса
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT id, name, points FROM routes"
	rows, err := r.db.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}

	// если в базе нет маршрутов
	if rows.Err() != nil {
		return nil, routes.ErrRouteNotFound
	}

	// в эту переменную будет сканиться результат запроса
	var id uuid.UUID
	var name string
	var points string

	// пробегаем по всем записям
	for rows.Next() {
		err = rows.Scan(&id, &name, &points)

		if err != nil {
			return nil, err
		}

		// создаем объект возвращаем его
		route, err := routes.NewRoute(id, name, points)

		if err != nil {
			return nil, err
		}

		arrRoutes = append(arrRoutes, *route)
	}

	return &arrRoutes, nil
}

// добавить маршрут
func (r *PostgresRepository) AddRoute(route *routes.Route) error {

	// Добавляем маршрут
	//**********************************
	// создаю контекст для запроса
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "INSERT INTO routes(id, name, points) VALUES ($1, $2, $3)"
	_, err := r.db.ExecContext(ctx, query, route.ID(), route.Name(), route.Points())
	if err != nil {
		// проверяем ошибку на предмет вставки маршрута с названием, которое уже есть в БД
		// создаем объект *pgconn.PgError - в нем будет храниться код ошибки из БД
		// QUESTION: если я вставляю запись с уже  существующим названием, то эту ошибку я получаю - т.к. название UNIQUE
		// но эта же самая ошибка вылезает если я вставляю запись с таким же uuid. Как различить эти ошибки?
		// т.е. как тут различить - ошибка из-за вставки с существующим именем или с существующим ID?
		var pgErr *pgconn.PgError

		// преобразуем ошибку к типу pgconn.PgError
		if errors.As(err, &pgErr) {
			// если ошибка- запись существует, то возвращаем эту ошибку
			if pgErr.Code == pgerrcode.UniqueViolation {
				return routes.ErrNameAlreadyExist
			}
		} else {
			return err
		}
	}

	return nil
}

// Сохранить маршрут
func (r *PostgresRepository) SaveRoute(route *routes.Route) error {

	// создаю контекст для запроса
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "UPDATE routes SET name=$1, points=$2 WHERE id=$3"
	_, err := r.db.ExecContext(ctx, query, route.Name(), route.Points(), route.ID())

	if err != nil {
		return err
	}

	return nil
}

// Удалить маршрут
func (r *PostgresRepository) DeleteRoute(id uuid.UUID) error {

	// создаю контекст для запроса
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "DELETE FROM routes WHERE id=$1"
	_, err := r.db.ExecContext(ctx, query, id)

	if err != nil {
		return err
	}

	return nil
}

// Найти маршрут по id
func (r *PostgresRepository) GetRouteByID(id uuid.UUID) (*routes.Route, error) {

	// создаю контекст для запроса
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT name, points FROM routes WHERE id=$1"
	row := r.db.QueryRowContext(ctx, query, id)

	// в эту переменную будет сканиться результат запроса
	var name string
	var points string

	err := row.Scan(&name, &points)

	if err != nil {
		return nil, err
	}

	// создаем маршрут и возвращаем его
	route, err := routes.NewRoute(id, name, points)

	if err != nil {
		return nil, err
	}

	return route, nil
}
