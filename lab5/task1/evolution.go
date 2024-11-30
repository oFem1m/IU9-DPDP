package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	rows       = 10
	cols       = 10
	steps      = 10000
	numThreads = 4
)

type Matrix [][]int

func createMatrix(rows, cols int) Matrix {
	matrix := make(Matrix, rows)
	for i := range matrix {
		matrix[i] = make([]int, cols)
		for j := range matrix[i] {
			matrix[i][j] = rand.Intn(2) // 0 или 1
		}
	}
	return matrix
}

// countNeighbors считает количество живых соседей для клетки (x, y)
func countNeighbors(matrix Matrix, x, y int) int {
	rows, cols := len(matrix), len(matrix[0])
	directions := [][2]int{{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 0}, {1, 1}}
	count := 0
	for _, d := range directions {
		nx, ny := (x+d[0]+rows)%rows, (y+d[1]+cols)%cols
		count += matrix[nx][ny]
	}
	return count
}

// evolvePart обрабатывает часть матрицы и вычисляет новое состояние клеток
func evolvePart(startRow, endRow int, matrix, newMatrix Matrix, barrier *sync.Cond, stepDone *int, totalThreads int, wg *sync.WaitGroup, stepTimes *[]time.Duration) {
	defer wg.Done()

	for step := 0; step < steps; step++ {
		startTime := time.Now()

		// Обрабатываем клетки
		for x := startRow; x < endRow; x++ {
			for y := 0; y < len(matrix[0]); y++ {
				neighbors := countNeighbors(matrix, x, y)
				oldState := matrix[x][y]
				newState := oldState

				// Применяем правила
				if oldState == 1 && (neighbors < 2 || neighbors > 3) {
					newState = 0
				} else if oldState == 0 && neighbors == 3 {
					newState = 1
				}

				newMatrix[x][y] = newState
			}
		}

		// Синхронизация с барьером
		barrier.L.Lock()
		*stepDone++
		if *stepDone == totalThreads {
			// Последний поток уведомляет всех
			*stepDone = 0
			for i := range matrix {
				copy(matrix[i], newMatrix[i])
			}
			(*stepTimes)[step] += time.Since(startTime) // Замер времени для текущего шага
			barrier.Broadcast()
		} else {
			// Остальные потоки ждут
			barrier.Wait()
		}
		barrier.L.Unlock()
	}
}

func printMatrix(matrix Matrix) {
	for _, row := range matrix {
		for _, cell := range row {
			fmt.Print(cell, " ")
		}
		fmt.Println()
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	matrix := createMatrix(rows, cols)
	newMatrix := createMatrix(rows, cols)

	fmt.Println("Initial state:")
	printMatrix(matrix)

	// Барьерная синхронизация
	barrier := sync.NewCond(&sync.Mutex{})
	stepDone := 0
	wg := &sync.WaitGroup{}
	stepTimes := make([]time.Duration, steps)

	// Запуск потоков
	chunkSize := rows / numThreads
	wg.Add(numThreads)
	for i := 0; i < numThreads; i++ {
		startRow := i * chunkSize
		endRow := startRow + chunkSize
		if i == numThreads-1 {
			endRow = rows // Последний поток берет оставшиеся строки
		}
		go evolvePart(startRow, endRow, matrix, newMatrix, barrier, &stepDone, numThreads, wg, &stepTimes)
	}

	// Ожидание завершения всех потоков
	wg.Wait()

	// Вычисление среднего времени шага
	var totalTime time.Duration
	for _, t := range stepTimes {
		totalTime += t
	}
	averageStepTime := totalTime / time.Duration(steps)

	fmt.Println("\nSimulation completed:")
	printMatrix(matrix)
	fmt.Printf("\nAverage step time: %v\n", averageStepTime)
}
