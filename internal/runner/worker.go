package runner

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
)

type Scheduler struct {
	tasksCh   chan Task
	wg        sync.WaitGroup
	completed uint64 // Счетчик успешно завершенных
	panics    uint64 // Счетчик паник
}

func NewScheduler(bufferSize int) *Scheduler {
	return &Scheduler{
		tasksCh: make(chan Task, bufferSize),
	}
}

// Добавь эти методы, чтобы main мог прочитать результат
func (s *Scheduler) GetStats() (completed, panics uint64) {
	return atomic.LoadUint64(&s.completed), atomic.LoadUint64(&s.panics)
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
					if !ok {
						return
					} // Канал закрыт
					func() {
						// Сначала определяем возврат в пул через defer
						// Он сработает ПОСЛЕ recover, когда выполнение функции завершится
						defer func() {
							if st, ok := task.(*SimpleTask); ok {
								// Очищаем данные перед возвратом (важно!)
								st.ID = ""
								TaskPool.Put(st)
							}
						}()
						defer func() {
							if r := recover(); r != nil {
								atomic.AddUint64(&s.panics, 1) // Атомарно +1
								fmt.Printf("[Worker %d] ПАНИКА поймана: %v, task: %v\n", workerID, r, task.GetID())
							}
						}()
						if err := task.Execute(ctx); err == nil {
							atomic.AddUint64(&s.completed, 1) // Атомарно +1
						}
					}()
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
