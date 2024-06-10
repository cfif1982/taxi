package routes

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

	query := fmt.Sprintf("INSERT INTO routes(id, name, points) VALUES (%v, '%v', '%v')", route.ID(), route.Name(), route.Points())
	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		// проверяем ошибку на предмет вставки маршрута с названием, которое уже есть в БД
		// QUESTION: или с ID, который тоже уже есть. Как тут различить - ошибка из-за вставки с существующим именем или с существующим ID?
		// создаем объект *pgconn.PgError - в нем будет храниться код ошибки из БД
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

// Редактировать маршрут
func (r *PostgresRepository) EditRoute(route *routes.Route) error {

	// создаю контекст для запроса
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := fmt.Sprintf("UPDATE routes "+
		"SET name='%v', points='%v' WHERE id=%v", route.Name(), route.Points(), route.ID())

	_, err := r.db.ExecContext(ctx, query)

	if err != nil {
		return err
	}

	return nil
}

// Удалить маршрут
// QUESTION: для удаления достаточно предать id удаляемого объекта? или нужно передавать всё-равно сам объект?
func (r *PostgresRepository) DeleteRoute(route *routes.Route) error {

	// создаю контекст для запроса
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := fmt.Sprintf("DELETE FROM routes WHERE id=%v", route.ID())

	_, err := r.db.ExecContext(ctx, query)

	if err != nil {
		return err
	}

	return nil
}

// Найти маршрут по id
func (r *PostgresRepository) GetRouteById(id uuid.UUID) (*routes.Route, error) {

	// создаю контекст для запроса
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := fmt.Sprintf("SELECT id, name, points FROM route WHERE id=%v", id)
	row := r.db.QueryRowContext(ctx, query)

	// в эту переменную будет сканиться результат запроса
	var name string
	var points string

	err := row.Scan(&name, &points)

	if err != nil {
		return nil, err
	}

	// создаем объект ссылку и возвращаем ее
	route, err := routes.NewRoute(id, name, points)

	if err != nil {
		return nil, err
	}

	return route, nil
}
