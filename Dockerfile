# BUILDING PROJECT
FROM golang:latest AS builder
WORKDIR /src
COPY . .
RUN GOOS=linux go build -o server ./cmd/main.go

# DOCKER STAGE: COPY NEEDED ELEMENTS TO NEW CONTAINER
FROM ubuntu:20.04
WORKDIR /src
COPY --from=builder ./src/ .
CMD ["./server"]