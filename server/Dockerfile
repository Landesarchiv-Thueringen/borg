FROM golang:alpine3.21 AS build
WORKDIR /borg
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o borg_server cmd/server.go

FROM alpine:3.21 AS prod
WORKDIR /borg
COPY --from=build /borg/borg_server ./borg_server
CMD ["./borg_server"]