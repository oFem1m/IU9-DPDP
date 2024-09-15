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

	var wg sync.WaitGroup

	rowsPerWorker := n / workers

	for w := 0; w < workers; w++ {
		startRow := w * rowsPerWorker
		endRow := startRow + rowsPerWorker
		if w == workers-1 {
			endRow = n
		}

		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for i := start; i < end; i++ {
				for j := 0; j < n; j++ {
					for k := 0; k < n; k++ {
						C[i][j] += A[i][k] * B[k][j]
					}
				}
			}
		}(startRow, endRow)
	}

	wg.Wait()

	return C
}

func main() {
	n := 500

	workers := 15

	A := createMatrix(n)
	B := createMatrix(n)

	start := time.Now()
	parallelMultiplyMatrices(A, B, n, workers)
	duration := time.Since(start)

	fmt.Printf("Время выполнения параллельного алгоритма (%d потоков): %v\n", workers, duration)
}
