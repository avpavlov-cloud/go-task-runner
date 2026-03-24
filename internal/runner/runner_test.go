package runner

import (
	"context"
	"testing"
	"time"
)

// MockTask для тестов, который ничего не делает, но имитирует работу
type MockTask struct{}

func (m *MockTask) Execute(ctx context.Context) error {
	return nil
}

// Добавьте этот метод:
func (m *MockTask) GetID() string {
    return "mock-task-id"
}

func TestScheduler(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	sched := NewScheduler(10)
	sched.Start(ctx, 5) // 5 воркеров

	taskCount := 100
	for i := 0; i < taskCount; i++ {
		sched.Submit(&MockTask{})
	}

	sched.Wait()

	completed, panics := sched.GetStats()
	if completed != uint64(taskCount) {
		t.Errorf("Ожидалось %d выполненных задач, получили %d", taskCount, completed)
	}
	if panics != 0 {
		t.Errorf("Ожидалось 0 паник, получили %d", panics)
	}
}

// Бенчмарк: проверяем, как быстро планировщик переваривает 1000 задач
func BenchmarkScheduler(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer() // Не считаем время подготовки
		ctx := context.Background()
		sched := NewScheduler(100)
		sched.Start(ctx, 10)
		b.StartTimer() // Начинаем замер

		for j := 0; j < 1000; j++ {
			sched.Submit(&MockTask{})
		}
		sched.Wait()
	}
}
