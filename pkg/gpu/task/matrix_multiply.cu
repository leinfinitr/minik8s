#include <stdio.h>
#include <stdlib.h>
#include <iostream>
#include <cuda_runtime.h>
#include <device_launch_parameters.h>

using namespace std;

const int M = 8;
const int N = 8;

__global__ void matrix_multiply(int *A, int *B, int *C, int width) {
    int row = blockIdx.y * blockDim.y + threadIdx.y;
    int col = blockIdx.x * blockDim.x + threadIdx.x;

    if (row < width && col < width) {
        int sum = 0;
        for (int k = 0; k < width; k++) {
            sum += A[row * width + k] * B[k * width + col];
        }
        C[row * width + col] = sum;
    }
}

int main() {
    int size = M * N * sizeof(int);

    int *host_A = (int *)malloc(size);
    int *host_B = (int *)malloc(size);
    int *host_C = (int *)malloc(size);

    for (int i = 0; i < M; i++) {
        for (int j = 0; j < N; j++) {
            host_A[i * N + j] = i * N + j;
            host_B[i * N + j] = i * N + j;
            host_C[i * N + j] = 0;
        }
    }

    // 打印初始矩阵
    cout << "矩阵A:" << endl;
    for (int i = 0; i < M; i++) {
        for (int j = 0; j < N; j++) {
            cout << host_A[i * N + j] << " ";
        }
        cout << endl;
    }

    cout << "矩阵B:" << endl;
    for (int i = 0; i < M; i++) {
        for (int j = 0; j < N; j++) {
            cout << host_B[i * N + j] << " ";
        }
        cout << endl;
    }

    // 分配device内存
    int *dev_A, *dev_B, *dev_C;
    cudaMalloc((void **)&dev_A, size);
    cudaMalloc((void **)&dev_B, size);
    cudaMalloc((void **)&dev_C, size);

    // 复制内存到device
    cudaMemcpy(dev_A, host_A, size, cudaMemcpyHostToDevice);
    cudaMemcpy(dev_B, host_B, size, cudaMemcpyHostToDevice);

    // 设置grid和block
    dim3 grid(M / 2, N / 2);
    dim3 block(2, 2);

    matrix_multiply<<<grid, block>>>(dev_A, dev_B, dev_C, N);

    cudaMemcpy(host_C, dev_C, size, cudaMemcpyDeviceToHost);

    cout << "乘法结果:" << endl;
    for (int i = 0; i < M; i++) {
        for (int j = 0; j < N; j++) {
            cout << host_C[i * N + j] << " ";
        }
        cout << endl;
    }

    // 释放内存
    free(host_A);
    free(host_B);
    free(host_C);
    cudaFree(dev_A);
    cudaFree(dev_B);
    cudaFree(dev_C);

    return 0;
}