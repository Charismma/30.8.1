package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Хранилище данных.
type Storage struct {
	db *pgxpool.Pool
}

// Конструктор, принимает строку подключения к БД.
func New(constr string) (*Storage, error) {
	db, err := pgxpool.Connect(context.Background(), constr)
	if err != nil {
		return nil, err
	}
	s := Storage{
		db: db,
	}
	return &s, nil
}

// Задача.
type Task struct {
	ID         int
	Opened     int64
	Closed     int64
	AuthorID   int
	AssignedID int
	Title      string
	Content    string
}

// Tasks возвращает список задач из БД. Всех задач,а также по автору
func (s *Storage) Tasks(taskID, authorID int) ([]Task, error) {
	rows, err := s.db.Query(context.Background(), `
		SELECT 
			id,
			opened,
			closed,
			author_id,
			assigned_id,
			title,
			content
		FROM tasks
		WHERE
			($1 = 0 OR id = $1) AND
			($2 = 0 OR author_id = $2)
		ORDER BY id;
	`,
		taskID,
		authorID,
	)
	if err != nil {
		return nil, err
	}
	var tasks []Task
	// итерирование по результату выполнения запроса
	// и сканирование каждой строки в переменную
	for rows.Next() {
		var t Task
		err = rows.Scan(
			&t.ID,
			&t.Opened,
			&t.Closed,
			&t.AuthorID,
			&t.AssignedID,
			&t.Title,
			&t.Content,
		)
		if err != nil {
			return nil, err
		}
		// добавление переменной в массив результатов
		tasks = append(tasks, t)

	}
	// ВАЖНО не забыть проверить rows.Err()
	return tasks, rows.Err()
}

// NewTask создаёт новую задачу и возвращает её id.
func (s *Storage) NewTask(t Task) (int, error) {
	var id int
	err := s.db.QueryRow(context.Background(), `
		INSERT INTO tasks (Opened,Closed,Author_ID,Assigned_ID,Title, Content)
		VALUES ($1, $2,$3,$4,$5,$6) RETURNING id;
		`,
		t.Opened,
		t.Closed,
		t.AuthorID,
		t.AssignedID,
		t.Title,
		t.Content,
	).Scan(&id)
	return id, err
}

// Получение списка задач по метке
func (s *Storage) TaskByLabelID(id int) ([]Task, error) {
	rows, err := s.db.Query(context.Background(),
		`SELECT tasks.* FROM tasks LEFT JOIN
		tasks_labels ON tasks_labels.task_id=tasks.id
		LEFT JOIN labels ON tasks_labels.label_id=labels.id
		WHERE labels.id=$1;`,
		id)
	if err != nil {
		return nil, err
	}
	var task []Task
	for rows.Next() {
		var t Task
		err = rows.Scan(
			&t.ID,
			&t.Opened,
			&t.Closed,
			&t.AuthorID,
			&t.AssignedID,
			&t.Title,
			&t.Content,
		)
		if err != nil {
			return nil, err
		}
		task = append(task, t)
	}

	return task, nil

}

// Обновление контента задачи по id
func (s *Storage) UpdateTaskContentById(id int, content string) (int, string, error) {
	var ids int
	var contents string
	err := s.db.QueryRow(context.Background(),
		`UPDATE tasks SET content=$1 WHERE id=$2 RETURNING id,content;`,
		content,
		id).Scan(&ids, &contents)

	return ids, contents, err

}

// Удаление задачи по ее id
func (s *Storage) DeleteTaskById(id int) error {
	_, err := s.db.Exec(context.Background(),
		`DELETE FROM tasks_labels WHERE task_id=$1;`,
		id)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(context.Background(),
		`DELETE FROM tasks WHERE id=$1;`,
		id)
	if err != nil {
		return err
	}
	fmt.Println("Удаление прошло успешно")
	return nil

}

// Вставка массива задач
func (s *Storage) NewSomeTasks(tasks []Task) ([]int, error) {
	var id int
	var ids []int
	var err error
	for _, task := range tasks {
		err = s.db.QueryRow(context.Background(), `
		INSERT INTO tasks (Opened,Closed,Author_ID,Assigned_ID,Title, Content)
		VALUES ($1, $2,$3,$4,$5,$6) RETURNING id;
		`,
			task.Opened,
			task.Closed,
			task.AuthorID,
			task.AssignedID,
			task.Title,
			task.Content,
		).Scan(&id)
		ids = append(ids, id)
	}

	return ids, err
}
