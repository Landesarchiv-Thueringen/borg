FROM golang:alpine3.20 AS build
WORKDIR /borg
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o borg cmd/server.go

FROM alpine:3.20 AS prod
WORKDIR /borg
COPY --from=build /borg/borg ./borg
CMD ["./borg"]