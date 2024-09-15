package main

import (
	"fmt"
	"math/rand"
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

func multiplyMatrices(A, B [][]int, n int) [][]int {
	C := make([][]int, n)
	for i := 0; i < n; i++ {
		C[i] = make([]int, n)
		for j := 0; j < n; j++ {
			for k := 0; k < n; k++ {
				C[i][j] += A[i][k] * B[k][j]
			}
		}
	}
	return C
}

func main() {
	n := 500

	A := createMatrix(n)
	B := createMatrix(n)

	start := time.Now()
	multiplyMatrices(A, B, n)
	duration := time.Since(start)

	// Выводим время выполнения
	fmt.Printf("Время выполнения стандартного алгоритма: %v\n", duration)
}
