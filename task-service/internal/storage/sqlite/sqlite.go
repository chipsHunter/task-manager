package sqlite

import (
	"database/sql"
	"fmt"
	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
	"task-service/internal/model/task"
	"task-service/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.new"
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS task(
    		id INTEGER PRIMARY KEY,
    		name TEXT NOT NULL,
		    description TEXT NOT NULL,
		    author TEXT NOT NULL,
		    type TEXT NOT NULL,
		    status TEXT NOT NULL
		);
	`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}
func (s *Storage) SaveTask(taskToSave task.Task) error {
	const op = "storage.sqlite.SaveTask"

	stmt, err := s.db.Prepare("INSERT INTO task(name, description,status,author,type) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(taskToSave.Name,
		taskToSave.Description,
		taskToSave.Status,
		taskToSave.Author,
		taskToSave.Type)
	if err != nil {
		//TODO : refactor this
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return fmt.Errorf("%s: %w", op, storage.ErrURLExists)
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) UpdateTask(taskToUpdate task.Task) error {
	const op = "storage.sqlite.UpdateTask"

	stmt, err := s.db.Prepare("UPDATE task SET (name, description, status, author, type) = (?,?,?,?,?) WHERE id = ?")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(
		taskToUpdate.Name,
		taskToUpdate.Description,
		taskToUpdate.Status,
		taskToUpdate.Author,
		taskToUpdate.Type,
		taskToUpdate.Id,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	numRows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if numRows == 0 {
		return fmt.Errorf("%s: task with id %d not found", op, taskToUpdate.Id)
	}

	return nil
}

func (s *Storage) DeleteTask(id int) error {
	const op = "storage.sqlite.DeleteTask"

	stmt, err := s.db.Prepare("DELETE FROM task WHERE id = ?")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	numRows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if numRows == 0 {
		return fmt.Errorf("%s: task with id %d not found", op, id)
	}

	return nil
}

func (s *Storage) GetTask(id int) (task.Task, error) {
	const op = "storage.sqlite.GetTask"

	var t task.Task

	row := s.db.QueryRow("SELECT id, name, description, status, author, type FROM task WHERE id = ?", id)
	if err := row.Scan(&t.Id, &t.Name, &t.Description, &t.Status, &t.Author, &t.Type); err != nil {
		if err == sql.ErrNoRows {
			return t, fmt.Errorf("%s: task with id %d not found", op, id)
		}
		return t, fmt.Errorf("%s: %w", op, err)
	}

	return t, nil
}

func (s *Storage) GetUserTasks(author string) ([]task.Task, error) {
	const op = "storage.sqlite.GetUserTasks"

	rows, err := s.db.Query("SELECT id, name, description, status, author, type FROM task WHERE author = ?", author)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var tasks []task.Task
	for rows.Next() {
		var t task.Task
		if err := rows.Scan(&t.Id, &t.Name, &t.Description, &t.Status, &t.Author, &t.Type); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return tasks, nil
}
