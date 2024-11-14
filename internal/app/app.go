package app

import (
	"TasksManager/internal/data_base"
	"context"
	"fmt"
	"log"
)

var ctx = context.Background()

func StartApp() {
	pool := db.InitDB("postgres://daniel_krerider:admin@127.0.0.1:5432/task?sslmode=disable")
	storage := &db.Storage{Db: pool}

	//	Пример добавления задачи
	if err := AddTask(storage); err != nil {
		log.Fatalf("Ошибка при добавлении задачи: %v", err)
	}

	//	Пример удаления задачи
	if err := DeleteTask(storage, 15); err != nil {
		log.Fatalf("Ошибка при удалении задачи: %v", err)
	}

	//	Пример обновления задачи
	if err := UpdateTask(storage); err != nil {
		log.Fatalf("Ошибка при обновлении задачи: %v", err)
	}

	//	Пример получения задачи с фильтрацией по метке и автору
	if err := GetFilteredTasks(storage); err != nil {
		log.Fatalf("Ошибка при получении задачи2: %v", err)
	}
}

// Добарление задачь
func AddTask(storage *db.Storage) error {
	task := db.Tasks{
		Title:      "Задача 3",
		Content:    "Текст задачи3",
		AuthorID:   1,
		AssignedID: 2,
	}

	id, err := storage.Newtask(ctx, task)
	if err != nil {
		return fmt.Errorf("не удалось добавить задачу: %v", err)
	}

	fmt.Printf("Задача успешно добавлена с ID: %d\n", id)
	return nil
}

// Удаление задачь по ID
func DeleteTask(storage *db.Storage, taskID int) error {
	taskToDelete := db.Tasks{ID: taskID}
	err := storage.DeleteTask(ctx, taskToDelete)
	if err != nil {
		return fmt.Errorf("не удалось удалить задачу: %w", err)
	}
	fmt.Printf("Задача с ID %d успешно удалена.\n", taskID)
	return nil
}

// Обновление задачи по ID
func UpdateTask(storage *db.Storage) error {

	taskToUpdate := db.Tasks{ID: 14, Title: "Обновленнная задача 2", Content: "купить бумагу а2 в офис"}

	err := storage.UpdateTask(ctx, taskToUpdate)
	if err != nil {
		return fmt.Errorf("Не удалось обновить значение %v\n", err)
	}
	fmt.Printf("Задача с ID %d успешно обновлена \n", taskToUpdate.ID)
	return nil
}

// Получение задачь
func GetFilteredTasks(storage *db.Storage) error {
	labelFilter := "Urgent" // Фильтр по метке (или '' для игнорирования)
	authorID := 1           // Фильтр по автору (или 0 для игнорирования)

	data, err := storage.GetTasks(labelFilter, authorID)
	if err != nil {
		return fmt.Errorf("Ошибка получения задач:%w ", err)
	}

	for _, task := range data {
		fmt.Printf("Задача ID: %d, Заголовок: %s, Автор: %s, Назначенный: %s\n",
			task.ID, task.Title, task.Author, task.Assigned)
	}
	return nil
}
