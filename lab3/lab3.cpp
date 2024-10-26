#include <iostream>
#include <vector>
#include <cmath>
#include <numeric>
#include <omp.h>

void parallel_simple_iteration(const std::vector<std::vector<double>>& A, const std::vector<double>& b,
                               double tau, double epsilon, int max_iter, std::vector<double>& x) {
    int N = b.size();
    std::vector<double> new_x(N, 0.0); // вектор для новых значений x

    for (int iter = 0; iter < max_iter; ++iter) {
        // шаг по всем элементам x
        #pragma omp parallel for
        for (int i = 0; i < N; ++i) {
            double row_sum = 0.0;
            // произведение строки из A на x
            for (int j = 0; j < N; ++j) {
                row_sum += A[i][j] * x[j];
            }
            // обновляем x[i]
            new_x[i] = x[i] - tau * (row_sum - b[i]);
        }

        // вычисляем норму
        double norm_diff = 0.0;
        #pragma omp parallel for reduction(+:norm_diff)
        for (int i = 0; i < N; ++i) {
            double diff = 0.0;
            for (int j = 0; j < N; ++j) {
                diff += A[i][j] * new_x[j];
            }
            norm_diff += (diff - b[i]) * (diff - b[i]);
        }
        norm_diff = std::sqrt(norm_diff) / std::sqrt(std::inner_product(b.begin(), b.end(), b.begin(), 0.0));

        // условие завершения
        if (norm_diff < epsilon) {
            std::cout << "Решение найдено за " << iter + 1 << " итераций" << std::endl;
            break;
        }

        // Обновляем x
        x = new_x;
    }
}

void test_known_solution(int N, double tau, double epsilon) {
    std::vector<std::vector<double>> A(N, std::vector<double>(N, 1.0));
    for (int i = 0; i < N; ++i) {
        A[i][i] = 2.0;
    }
    std::vector<double> b(N, N + 1);
    std::vector<double> x(N, 0.0);

    double start_time = omp_get_wtime();
    parallel_simple_iteration(A, b, tau, epsilon, 10000, x);
    double end_time = omp_get_wtime();

    std::cout << "Задача с известным решением:\n";
    std::cout << "Решение: ";
    for (double xi : x) {
        std::cout << xi << " ";
    }
    std::cout << "\nВремя выполнения теста: " << end_time - start_time << " секунд" << std::endl;
}

void test_arbitrary_solution(int N, double tau, double epsilon) {
    std::vector<std::vector<double>> A(N, std::vector<double>(N, 1.0));
    for (int i = 0; i < N; ++i) {
        A[i][i] = 2.0;
    }
    std::vector<double> u(N);
    for (int i = 0; i < N; ++i) {
        u[i] = sin(2 * M_PI * i / N);
    }

    std::vector<double> b(N, 0.0);
    for (int i = 0; i < N; ++i) {
        for (int j = 0; j < N; ++j) {
            b[i] += A[i][j] * u[j];
        }
    }
    std::vector<double> x(N, 0.0);

    double start_time = omp_get_wtime();
    parallel_simple_iteration(A, b, tau, epsilon, 10000, x);
    double end_time = omp_get_wtime();

    std::cout << "Задача с произвольным решением:\n";
    std::cout << "Решение: ";
    for (double xi : x) {
        std::cout << xi << " ";
    }
    std::cout << "\nВремя выполнения теста: " << end_time - start_time << " секунд" << std::endl;
}

int main() {
    int num_threads = 4;
    omp_set_num_threads(num_threads);

    std::cout << "Используется потоков: " << num_threads << std::endl;

    int N = 100;
    double epsilon = 0.0001;
    double tau = 0.1 / N;

    // тесты
    test_known_solution(N, tau, epsilon);
    test_arbitrary_solution(N, tau, epsilon);

    return 0;
}
