package main

import (
	"context"
	"fmt"
	"taskrunner/internal/runner"
	"time"
)

func main() {
	// Контекст закроется сам через 3 секунды
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	sched := runner.NewScheduler(10)

	// Запускаем 3 воркера
	sched.Start(ctx, 3)

	// Накидываем 5 задач
	for i := 1; i <= 5; i++ {
		sched.Submit(&runner.SimpleTask{ID: fmt.Sprintf("%d", i)})
	}

	sched.Wait()
	c, p := sched.GetStats()
	fmt.Printf("Итог: Выполнено: %d, Паник: %d\n", c, p)
}
