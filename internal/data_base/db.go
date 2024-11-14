package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

var ctx context.Context = context.Background()

func InitDB(dataSourceName string) *pgxpool.Pool {
	db, err := pgxpool.Connect(ctx, dataSourceName)
	if err != nil {
		log.Fatalf("Ошибка подключения %v\n", err)
	}
	fmt.Println("Успешное подключение!")
	return db // Возвращаем пул соединений
}

type Tasks struct {
	ID         int
	Opened     int
	Closed     int
	Title      string
	Content    string
	AuthorID   int
	AssignedID int
	Author     string
	Assigned   string
}

type Storage struct {
	Db *pgxpool.Pool
}

// Добавление новой задачи
func (s *Storage) Newtask(ctx context.Context, t Tasks) (int, error) {
	var id int

	err := s.Db.QueryRow(ctx,
		`INSERT INTO tasks (title, content, author_id, assigned_id) VALUES ($1, $2, $3, $4) RETURNING id;`,
		t.Title, t.Content, t.AuthorID, t.AssignedID).
		Scan(&id)

	return id, err
}

// Удаление задачи по ID
func (s *Storage) DeleteTask(ctx context.Context, t Tasks) error {
	result, err := s.Db.Exec(ctx, `DELETE FROM tasks WHERE id=$1;`, t.ID)
	if err != nil {
		log.Fatal(err)
	}
	rowsAffective := result.RowsAffected()
	if rowsAffective == 0 {
		return fmt.Errorf("задача с ID %d не найдена", t.ID)
	}
	return nil
}

// Обновление задачи по ID
func (s *Storage) UpdateTask(ctx context.Context, t Tasks) error {
	result, err := s.Db.Exec(ctx, `UPDATE tasks SET title=$1,content=$2 WHERE id=$3;`, t.Title, t.Content, t.ID)
	if err != nil {
		log.Fatal(err)
	}
	rowsAffective := result.RowsAffected()
	if rowsAffective == 0 {
		return fmt.Errorf("задача с ID %d не найдена", t.ID)
	}
	return nil
}

// Получение задач с фильтрацией по метке и автору
func (s *Storage) GetTasks(label string, authorID int) ([]Tasks, error) {
	rows, err := s.Db.Query(context.Background(), `
SELECT 
    tasks.id,
    tasks.opened,
    tasks.closed,
    tasks.title,
    tasks.content, 
    author.name AS author,
    assigned.name AS assigned
FROM tasks
JOIN users AS author ON author_id = author.id
JOIN users AS assigned ON assigned_id = assigned.id
JOIN tasks_labels ON tasks.id = tasks_labels.task_id
JOIN labels ON tasks_labels.label_id = labels.id
		WHERE ($1 = '' OR labels.label = $1) AND ($2 = 0 OR tasks.author_id = $2)
		ORDER BY tasks.id;
	`, label, authorID)

	if err != nil {
		return nil, fmt.Errorf("не удалось выполнить запрос: %w", err)
	}
	defer rows.Close()

	var tasks []Tasks

	for rows.Next() {
		var t Tasks
		err = rows.Scan(
			&t.ID,
			&t.Opened,
			&t.Closed,
			&t.Title,
			&t.Content,
			&t.Author,
			&t.Assigned,
		)
		if err != nil {
			return nil, fmt.Errorf("не удалось сканировать строку: %w", err)
		}
		tasks = append(tasks, t)
	}

	return tasks, rows.Err()
}
