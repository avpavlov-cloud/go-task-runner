package main

import (
	"context"
	"fmt"
	"taskrunner/internal/runner"
)

func main() {
	ctx := context.Background()
	sched := runner.NewScheduler(10)

	// Запускаем 3 воркера
	sched.Start(ctx, 3)

	// Накидываем 5 задач
	for i := 1; i <= 5; i++ {
		sched.Submit(&runner.SimpleTask{ID: fmt.Sprintf("%d", i)})
	}

	sched.Wait()
	fmt.Println("Все задачи выполнены!")
}
