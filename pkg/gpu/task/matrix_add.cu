#include <stdio.h>
#include <stdlib.h>
#include <iostream>
#include <cuda_runtime.h>
#include <device_launch_parameters.h>

using namespace std;

#define M 7
#define N 7

// 矩阵加法 Kernel
__global__ void addMatrices(int *A, int *B, int *C, int m, int n) {
    int idx = threadIdx.x + blockIdx.x * blockDim.x;
    int idy = threadIdx.y + blockIdx.y * blockDim.y;
    if (idx < m && idy < n) {
        int index = idx * n + idy;
        C[index] = A[index] + B[index];
    }
}

// 初始化矩阵数据
void matrixInit(int* mat, int m, int n) {
    for(int i = 0; i < m; i++) {
        for(int j = 0; j < n; j++) {
            mat[i*n + j] = i*n+j;
        }
    }
}

void printMatrix(int *mat, int m, int n) {
    for(int i = 0; i < m; i++) {
        for(int j = 0; j < n; j++) {
            cout << mat[i*n + j] << "\t";
        }
        cout << endl;
    }
}

int main() {
    int size = M*N*sizeof(int);

    // 分配host内存
    int *host_A = (int *) malloc(size);
    int *host_B = (int *) malloc(size);
    int *host_C = (int *) malloc(size);

    // 初始化矩阵
    matrixInit(host_A, M, N);
    matrixInit(host_B, M, N);

    // 打印初始矩阵
    cout<<"矩阵A:"<<endl;
    printMatrix(host_A, M, N);
    cout<<"矩阵B:"<<endl;
    printMatrix(host_B, M, N);

    // 分配device内存
    int *dev_A, *dev_B, *dev_C;
    cudaMalloc((void **)&dev_A, size);
    cudaMalloc((void **)&dev_B, size);
    cudaMalloc((void **)&dev_C, size);

    // 复制内存到device
    cudaMemcpy((void *)dev_A, (void *)host_A, size, cudaMemcpyHostToDevice);
    cudaMemcpy((void *)dev_B, (void *)host_B, size, cudaMemcpyHostToDevice);

    // 设置grid和block
    dim3 grid((M+1)/2, (N+1)/2);
    dim3 block(2, 2);

    // Launch the kernel
    addMatrices<<<grid, block>>>(dev_A, dev_B, dev_C, M, N);

    // 复制结果回host 
    cudaMemcpy((void *)host_C, (void *)dev_C, size, cudaMemcpyDeviceToHost);

    // 打印结果
    cout<<"矩阵加法结果:"<<endl;
    printMatrix(host_C, M, N);

    // 释放内存
    free(host_A);
    free(host_B);
    free(host_C);
    cudaFree(dev_A);
    cudaFree(dev_B);
    cudaFree(dev_C);

    return 0;
}