package main

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof" // Магия: регистрирует обработчики pprof
	"taskrunner/internal/runner"
	"time"
)

func main() {
	// Запускаем pprof сервер в фоне
	go func() {
		fmt.Println("Pprof доступен на: http://localhost:6060/debug/pprof/")
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			fmt.Printf("Ошибка pprof сервера: %v\n", err)
		}
	}()

	// Контекст закроется сам через 3 секунды
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second) //Увеличил для профилировщика с 3 до 300
	defer cancel()

	sched := runner.NewScheduler(10)

	// Запускаем 3 воркера
	sched.Start(ctx, 3)

	// Накидываем 5 задач
	for i := 1; i <= 1000; i++ {
		// Берем объект из пула (утверждаем тип через .(type))
		task := runner.TaskPool.Get().(*runner.SimpleTask)

		// 2. Инициализируем новыми данными
		task.ID = fmt.Sprintf("%d", i)

		sched.Submit(task)
	}

	sched.Wait()
	c, p := sched.GetStats()
	fmt.Printf("Итог: Выполнено: %d, Паник: %d\n", c, p)
}
