package storage

import (
	"log"
	"os"
	"strconv"
	"testing"
)

var s *Storage

func TestMain(m *testing.M) {
	constr := "postgres://postgres:password@192.168.1.191:5432/tasks"
	var err error
	s, err = New(constr)
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(m.Run())
}
func TestStorage_Tasks(t *testing.T) {
	data, err := s.Tasks(0, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(data)
	data, err = s.Tasks(1, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(data)
}

func TestStorage_NewTask(t *testing.T) {
	nTask := Task{Opened: 0,
		Closed:     0,
		AuthorID:   0,
		AssignedID: 0,
		Title:      "Новая задача при запуске",
		Content:    "Содержание новой задачи"}
	id_newTask, err := s.NewTask(nTask)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Создана задача с id:", id_newTask)
	tasks, err := s.Tasks(0, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tasks)
}

func TestStorage_TaskByLabelID(t *testing.T) {
	n, err := s.TaskByLabelID(1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(n)
	n, err = s.TaskByLabelID(0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(n)

}

func TestStorage_UpdateTaskContentById(t *testing.T) {
	a, b, err := s.UpdateTaskContentById(1, "Приколюха")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("В задаче с id:", a, " Обновлен контент на: ", b)
}

func TestStorage_DeleteTaskById(t *testing.T) {
	err := s.DeleteTaskById(7)
	if err != nil {
		t.Fatal(err)
	}
}

func TestStorage_NewSomeTasks(t *testing.T) {
	var massive []Task
	for i := 1; i <= 10; i++ {
		massive = append(massive, Task{Opened: int64(i),
			Closed:     int64(i),
			AuthorID:   0,
			AssignedID: 0,
			Title:      strconv.Itoa(i),
			Content:    strconv.Itoa(i)})
	}
	ids, err := s.NewSomeTasks(massive)
	if err != nil {
		t.Fatal(err)
	}
	for _, id := range ids {
		t.Log("Вставлена запись и ее id: ", id)
	}
}
