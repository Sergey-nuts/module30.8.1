package storage

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
}

func New(config string) (*Storage, error) {
	db, err := pgxpool.Connect(context.Background(), config)
	if err != nil {
		return nil, err
	}
	s := Storage{db}
	return &s, nil
}

type Task struct {
	ID         int64
	Opened     int64
	Closed     int64
	AuthorID   int64
	AssignedID int64
	Title      string
	Content    string
}

// возвращает список задач
// можно искать по id задачи TaskID или по id автора AuthorID
func (s *Storage) Tasks(TaskID, AuthorID int64) ([]Task, error) {
	rows, err := s.db.Query(context.Background(), `
	SELECT id, opened, closed, author_id, assigned_id, title, content
	FROM tasks
	WHERE
		($1 = 0 OR $1 = id) AND
		($2 = 0 OR $2 = author_id)
	ORDER BY id;
	`,
		TaskID,
		AuthorID,
	)
	if err != nil {
		return nil, err
	}
	var tasks []Task
	var t Task
	for rows.Next() {
		err := rows.Scan(&t.ID, &t.Opened, &t.Closed, &t.AuthorID, &t.AssignedID, &t.Title, &t.Content)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return tasks, rows.Err() // не забывать проверить rows.Err()
}

// создает новую залачу, возврящает её id
func (s *Storage) NewTask(t Task) (int64, error) {
	var id int64
	err := s.db.QueryRow(context.Background(), `
		INSERT INTO tasks (title, content)
		VALUES ($1, $2)
		RETURNING id;	
		`, t.Title, t.Content,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// обновить задачу по id (все поля кроме id, opened, author_id)
func (s *Storage) UpdateTask(id int64, t Task) error {
	_, err := s.db.Exec(context.Background(), `
		UPDATE tasks
		SET 
			closed = $2,
			assigned_id = $3,
			title = $4,
			content = $5
		WHERE id = $1;
	`, id, t.Closed, t.AssignedID, t.Title, t.Content)
	if err != nil {
		return err
	}
	return nil
}

// удалить задачу по id
func (s *Storage) DeletTask(id int64) error {
	_, err := s.db.Exec(context.Background(), `
		DELETE FROM tasks
		WHERE id = $1	
	`, id)
	return err
}
