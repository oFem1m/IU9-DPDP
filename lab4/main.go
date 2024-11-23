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
	fmt.Printf("Философ %d: %s\n", p.id, p.state)
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
			p.act("думает", time.Duration(rand.Intn(1000))*time.Millisecond)

			// Взять левую вилку
			p.leftFork.Lock()
			p.act("взял левую вилку", time.Duration(rand.Intn(500))*time.Millisecond)

			// Попробовать взять правую вилку
			if !p.tryTakeRightFork() {
				p.leftFork.Unlock()
				p.act("не смог взять правую вилку, положил левую", time.Duration(rand.Intn(500))*time.Millisecond)
				continue
			}

			// Трапеза
			p.act("ест", time.Duration(rand.Intn(1000))*time.Millisecond)

			// Положить обе вилки
			p.rightFork.Unlock()
			p.leftFork.Unlock()
			p.act("положил вилки", time.Duration(rand.Intn(500))*time.Millisecond)
		}
	}
}

// Попытка взять правую вилку
func (p *Philosopher) tryTakeRightFork() bool {
	locked := make(chan bool, 1)
	go func() {
		p.rightFork.Lock()
		locked <- true
	}()
	select {
	case success := <-locked:
		return success
	case <-time.After(50 * time.Millisecond): // Таймаут для ожидания правой вилки
		return false
	}
}

func main() {
	const numPhilosophers = 5
	const simulationTime = 10 * time.Second

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

	fmt.Println("Надумались и наелись")
}
