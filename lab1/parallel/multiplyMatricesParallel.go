package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func createMatrix(n int) [][]int {
	matrix := make([][]int, n)
	for i := 0; i < n; i++ {
		matrix[i] = make([]int, n)
		for j := 0; j < n; j++ {
			matrix[i][j] = rand.Intn(10)
		}
	}
	return matrix
}

func parallelMultiplyMatrices(A, B [][]int, n, workers int) [][]int {
	C := make([][]int, n)
	for i := 0; i < n; i++ {
		C[i] = make([]int, n)
	}

	var wg sync.WaitGroup // для ожидания завершения всех горутин

	rowsPerWorker := n / workers // Количество строк, обрабатываемых каждой горутиной

	for w := 0; w < workers; w++ {
		startRow := w * rowsPerWorker      // Начальная строка для горутины
		endRow := startRow + rowsPerWorker // Конечная строка для горутины
		if w == workers-1 {
			endRow = n // остальные строки
		}

		wg.Add(1) // inc WaitGroup
		go func(start, end int) {
			defer wg.Done() // dec WaitGroup после завершения горутины
			for i := start; i < end; i++ {
				for j := 0; j < n; j++ {
					for k := 0; k < n; k++ {
						C[i][j] += A[i][k] * B[k][j]
					}
				}
			}
		}(startRow, endRow)
	}

	wg.Wait() // Ожидаем завершения всех горутин

	return C
}

func main() {
	n := 500

	workers := 7

	A := createMatrix(n)
	B := createMatrix(n)

	start := time.Now()
	parallelMultiplyMatrices(A, B, n, workers)
	duration := time.Since(start)

	fmt.Printf("Время выполнения параллельного алгоритма (%d потоков): %v\n", workers, duration)
}
