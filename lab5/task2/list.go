package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Node struct {
	value int
	next  *Node
}

type LinkedList struct {
	head *Node
	lock sync.RWMutex
}

// Contains проверяет, содержится ли заданное значение в списке
func (ll *LinkedList) Contains(value int) bool {
	for current := ll.head; current != nil; current = current.next {
		if current.value == value {
			return true
		}
	}
	return false
}

// AddIfNotExists добавляет значение в список, если его там еще нет
func (ll *LinkedList) AddIfNotExists(value int) {
	ll.lock.RLock() // Блокировка на чтение
	alreadyExists := ll.Contains(value)
	ll.lock.RUnlock()

	if alreadyExists {
		return
	}

	ll.lock.Lock() // Блокировка на запись
	defer ll.lock.Unlock()

	// Повторная проверка на случай, если значение было добавлено другим потоком
	if !ll.Contains(value) {
		newNode := &Node{value: value}
		newNode.next = ll.head
		ll.head = newNode
	}
}

// Print печатает список
func (ll *LinkedList) Print() {
	ll.lock.RLock()
	defer ll.lock.RUnlock()
	for current := ll.head; current != nil; current = current.next {
		fmt.Print(current.value, " -> ")
	}
	fmt.Println("nil")
}

// Validate проверяет, что нет повторяющихся значений
func (ll *LinkedList) Validate() bool {
	ll.lock.RLock()
	defer ll.lock.RUnlock()
	values := make(map[int]bool)
	for current := ll.head; current != nil; current = current.next {
		if values[current.value] {
			return false
		}
		values[current.value] = true
	}
	return true
}

// Worker генерирует случайные числа и добавляет их в список
func Worker(id int, ll *LinkedList, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 10; i++ {
		value := rand.Intn(1000) // Генерация случайного числа
		ll.AddIfNotExists(value) // Попытка добавить число в список
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	ll := &LinkedList{}
	var wg sync.WaitGroup
	numThreads := 4 // Количество потоков

	wg.Add(numThreads)
	for i := 0; i < numThreads; i++ {
		go Worker(i, ll, &wg) // Запуск потоков
	}
	wg.Wait() // Ожидание завершения всех потоков

	// Печать списка
	fmt.Println("Linked list contents:")
	ll.Print()

	// Проверка на дубликаты
	if ll.Validate() {
		fmt.Println("No duplicates found!")
	} else {
		fmt.Println("Duplicates found!")
	}
}
