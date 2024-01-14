# The golang image is used only to build the server
FROM golang AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -v -o ./server ./cmd/server/

# Copy the server file and run it inside a new container
FROM ubuntu
WORKDIR /
COPY ./assets ./assets
COPY .env .env
COPY --from=builder /app/server ./server
CMD ["./server"]