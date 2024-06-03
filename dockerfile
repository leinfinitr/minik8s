FROM golang:1.22.2 as builder
WORKDIR /minik8s
COPY . .
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go build -o /minik8s/pkg/gpu/task_server /minik8s/pkg/gpu/main
RUN cp /minik8s/pkg/gpu/task_server /bin/server

FROM ubuntu:22.04
COPY --from=builder /bin/server /bin/server
RUN echo "nameserver 223.5.5.5" >  /etc/resolv.conf
ENTRYPOINT [ "/bin/server" ]