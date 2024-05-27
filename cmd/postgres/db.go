package main

import (
	"fmt"
	"log"
	"strconv"

	"DB_30.8.1/pkg/storage"
)

func main() {
	//Подключаемся к базе(Создаем пул соединений)
	constr := "postgres://postgres:password@192.168.1.191:5432/tasks"
	db, err := storage.New(constr)
	if err != nil {
		log.Println(err)
		return
	}
	//Создаем новую задачу и выводим ее id
	nTask := storage.Task{Opened: 0,
		Closed:     0,
		AuthorID:   0,
		AssignedID: 0,
		Title:      "Новая задача при запуске",
		Content:    "Содержание новой задачи"}
	id_newTask, err := db.NewTask(nTask)
	if err != nil {
		log.Println("Ошибка при создании новой задачи: ", err)
		return
	}
	fmt.Println("id новой задачи", id_newTask)
	//Получение задач по id, и id автора. В данном случае вывод всех задач
	l, err := db.Tasks(0, 0)
	if err != nil {
		log.Println("Ошибка при выводе всех задач", err)
		return
	}
	fmt.Println("Вывод всех задач: ", l)
	//Вывод информации о задачах по id метки
	n, err := db.TaskByLabelID(1)
	if err != nil {
		log.Println("Ошибка при выводе конкретной задачи по id", err)
		return
	}
	fmt.Println("Вывод значения по id метки", n)
	//Обновление контента в задаче по id
	a, b, err := db.UpdateTaskContentById(1, "Приколюха")
	if err != nil {
		log.Println("Ошибка обновления задачи", err)
		return
	}
	fmt.Println("В задаче с id:", a, " Обновлен контент на: ", b)
	//Удаление задачи по id
	err = db.DeleteTaskById(4)
	if err != nil {
		log.Println("Ошибка удаления записи", err)
		return
	}
	//Создание массива задач
	var massive []storage.Task
	for i := 1; i <= 10; i++ {
		massive = append(massive, storage.Task{Opened: int64(i),
			Closed:     int64(i),
			AuthorID:   0,
			AssignedID: 0,
			Title:      strconv.Itoa(i),
			Content:    strconv.Itoa(i)})
	}
	ids, err := db.NewSomeTasks(massive)
	if err != nil {
		log.Println("Ошибка при вставке массива задач", err)
		return
	}
	for _, id := range ids {
		fmt.Println("Вставлена запись и ее id: ", id)
	}

}
