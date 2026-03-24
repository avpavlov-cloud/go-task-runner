package runner

import (
	"context"
	"fmt"
	"time"
)

// Task — интерфейс для любой задачи в системе
type Task interface {
	Execute(ctx context.Context) error
	GetID() string
}

// SimpleTask — простая реализация задачи для теста
type SimpleTask struct {
	ID string
}

func (s *SimpleTask) Execute(ctx context.Context) error {
	start := time.Now() // Засекаем время СТАРТА

	defer func() {
		duration := time.Since(start) // Считаем, сколько прошло
		fmt.Printf("[Task %s] Суммарное время в Execute: %v\n", s.ID, duration)
	}()

	fmt.Printf("[Task %s] Начинаю выполнение...\n", s.ID)
	time.Sleep(500 * time.Millisecond) // Имитация работы
	if s.ID == "3" {
		panic("Вызов паники!")
	}

	select {
	case <-time.After(2 * time.Second): // Имитация долгой работы
		fmt.Printf("[Task %s] Готово.\n", s.ID)
		return nil
	case <-ctx.Done(): // Если контекст отменили (тайм-аут)
		fmt.Printf("[Task %s] ПРЕРВАНО по таймауту!\n", s.ID)
		return ctx.Err()
	}
}

func (s *SimpleTask) GetID() string {
	return s.ID
}
