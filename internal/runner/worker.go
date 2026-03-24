package runner

import (
	"context"
	"sync"
)

type Scheduler struct {
	tasksCh chan Task
	wg      sync.WaitGroup
}

func NewScheduler(bufferSize int) *Scheduler {
	return &Scheduler{
		tasksCh: make(chan Task, bufferSize),
	}
}

// Start запускает N воркеров (горутин)
func (s *Scheduler) Start(ctx context.Context, workerCount int) {
	for i := 0; i < workerCount; i++ {
		s.wg.Add(1)
		go func(workerID int) {
			defer s.wg.Done()
			for {
				select {
				case <-ctx.Done(): // Остановка по контексту
					return
				case task, ok := <-s.tasksCh:
					if !ok { return } // Канал закрыт
					_ = task.Execute(ctx)
				}
			}
		}(i)
	}
}

// Submit отправляет задачу в канал
func (s *Scheduler) Submit(t Task) {
	s.tasksCh <- t
}

func (s *Scheduler) Wait() {
	close(s.tasksCh)
	s.wg.Wait()
}
