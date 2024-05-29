FROM golang:1.22.3 as builder
WORKDIR /minik8s
COPY . /minik8s/
RUN go build -o /minik8s/pkg/gpu/server /minik8s/cmd/server/main
RUN mv /minik8s/pkg/gpu/server /bin/server

FROM ubuntu:20.04
COPY --from=builder /bin/server /bin/server

ENTRYPOINT [ "/bin/server" ]