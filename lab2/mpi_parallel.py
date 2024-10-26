from mpi4py import MPI
import numpy as np
import time


def parallel_simple_iteration(A, b, tau, epsilon, max_iter=10000):
    comm = MPI.COMM_WORLD
    rank = comm.Get_rank()  # номер текущего процесса
    size = comm.Get_size()  # количество процессов

    N = len(b)
    # размер локальной части для каждого процесса
    local_N = N // size + (1 if rank < N % size else 0)
    local_start = sum(N // size + (1 if i < N % size else 0) for i in range(rank))
    local_end = local_start + local_N

    # делим матрицу по строкам и создаем локальные части
    local_A = A[local_start:local_end]
    local_b = b[local_start:local_end]
    local_x = np.zeros(local_N)

    # хранит итоговое значение x на всех процессах
    global_x = np.zeros(N)

    for n in range(max_iter):
        # рассылаем данные global_x на все процессы
        comm.Bcast(global_x, root=0)

        # локально вычисляем Ax
        local_Ax = np.dot(local_A, global_x)

        # обновляем локальные значения вектора x
        local_x_new = local_x - tau * (local_Ax - local_b)

        # собираем обновленные x из всех процессов
        comm.Allgather(local_x_new, global_x)

        # вычисляем норму разности
        norm_diff = np.linalg.norm(np.dot(A, global_x) - b) / np.linalg.norm(b)
        # распространяем норму на все процессы
        norm_diff = comm.bcast(norm_diff, root=0)

        # условие завершения
        if norm_diff < epsilon:
            if rank == 0:
                print(f"Решение найдено за {n + 1} итераций")
            break

        # обновляем локальный вектор x
        local_x = local_x_new

    return global_x if rank == 0 else None


def test_known_solution(N, tau, epsilon):
    A = np.full((N, N), 1.0)
    np.fill_diagonal(A, 2.0)
    b = np.full(N, N + 1)

    start_time = time.time()
    x_solution = parallel_simple_iteration(A, b, tau, epsilon)
    end_time = time.time()

    if MPI.COMM_WORLD.Get_rank() == 0:
        print("Модельная задача с известным решением:")
        print("Решение:", x_solution)
        print(f"Общее время выполнения: {end_time - start_time:.4f} секунд")


def test_arbitrary_solution(N, tau, epsilon):
    A = np.full((N, N), 1.0)
    np.fill_diagonal(A, 2.0)
    u = np.sin(2 * np.pi * np.arange(N) / N)
    b = A @ u  # Формируем b

    start_time = time.time()
    x_solution = parallel_simple_iteration(A, b, tau, epsilon)
    end_time = time.time()

    if MPI.COMM_WORLD.Get_rank() == 0:
        print("Модельная задача с произвольным решением:")
        print("Решение:", x_solution)
        print(f"Общее время выполнения: {end_time - start_time:.4f} секунд")


if __name__ == "__main__":
    total_start_time = time.time()

    N = 100
    epsilon = 0.0001
    tau = 0.1 / N

    # тесты
    test_known_solution(N, tau, epsilon)
    test_arbitrary_solution(N, tau, epsilon)

    total_end_time = time.time()
    total_elapsed_time = total_end_time - total_start_time

    if MPI.COMM_WORLD.Get_rank() == 0:
        print("Общее время выполнения всех тестов: {:.4f} секунд".format(total_elapsed_time))
