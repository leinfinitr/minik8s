#!/bin/bash

for i in {1..15}; do
    go run minik8s/pkg/kubectl/main serverless run add 5 &
done