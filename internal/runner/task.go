package runner

import (
	"context"
	"fmt"
	"time"
)

// Task — интерфейс для любой задачи в системе
type Task interface {
	Execute(ctx context.Context) error
}

// SimpleTask — простая реализация задачи для теста
type SimpleTask struct {
	ID string
}

func (s *SimpleTask) Execute(ctx context.Context) error {
	fmt.Printf("[Task %s] Начинаю выполнение...\n", s.ID)
	time.Sleep(500 * time.Millisecond) // Имитация работы
	if s.ID == "3" {
		panic("Вызов паники!")
	}
	fmt.Printf("[Task %s] Завершено.\n", s.ID)
	return nil
}
