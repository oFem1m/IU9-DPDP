package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Philosopher Философ
type Philosopher struct {
	id        int
	leftFork  *sync.Mutex
	rightFork *sync.Mutex
	state     string
	logMutex  *sync.Mutex
}

// Действие философа
func (p *Philosopher) act(action string, duration time.Duration) {
	p.state = action
	p.logState()
	time.Sleep(duration)
}

// Состояние философа
func (p *Philosopher) logState() {
	p.logMutex.Lock()
	defer p.logMutex.Unlock()
	fmt.Printf("Person %d: %s\n", p.id, p.state)
}

// Жизненный цикл философа
func (p *Philosopher) dine(done *sync.WaitGroup, stop <-chan bool) {
	defer done.Done()
	for {
		select {
		case <-stop:
			return
		default:
			// Размышление
			p.act("thinking", time.Duration(rand.Intn(1000))*time.Millisecond)

			if p.id%2 == 0 {
				// Четные философы сначала берут правую вилку
				p.rightFork.Lock()
				p.act("get right fork", time.Duration(rand.Intn(500))*time.Millisecond)
				p.leftFork.Lock()
				p.act("get left fork", time.Duration(rand.Intn(500))*time.Millisecond)
			} else {
				// Нечетные философы сначала берут левую вилку
				p.leftFork.Lock()
				p.act("get left fork", time.Duration(rand.Intn(500))*time.Millisecond)
				p.rightFork.Lock()
				p.act("get right fork", time.Duration(rand.Intn(500))*time.Millisecond)
			}

			// Трапеза
			p.act("meal", time.Duration(rand.Intn(1000))*time.Millisecond)

			// Положить вилки
			p.rightFork.Unlock()
			p.leftFork.Unlock()
			p.act("put forks", time.Duration(rand.Intn(500))*time.Millisecond)
		}
	}
}

func main() {
	const numPhilosophers = 5
	const simulationTime = 2 * time.Second

	// Создание вилок
	forks := make([]*sync.Mutex, numPhilosophers)
	for i := 0; i < numPhilosophers; i++ {
		forks[i] = &sync.Mutex{}
	}

	// Логи состояний
	logMutex := &sync.Mutex{}

	// Создание философов
	philosophers := make([]*Philosopher, numPhilosophers)
	for i := 0; i < numPhilosophers; i++ {
		philosophers[i] = &Philosopher{
			id:        i + 1,
			leftFork:  forks[i],
			rightFork: forks[(i+1)%numPhilosophers],
			logMutex:  logMutex,
		}
	}

	// Канал для остановки
	stop := make(chan bool)
	var wg sync.WaitGroup

	// Запуск философов
	for _, philosopher := range philosophers {
		wg.Add(1)
		go philosopher.dine(&wg, stop)
	}

	// Завершение
	time.Sleep(simulationTime)
	close(stop)
	wg.Wait()

	fmt.Println("finish")
}
