FROM golang:alpine3.20 AS build
WORKDIR /borg/tools/tika
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o tika_api cmd/tika_api.go

FROM alpine:3.20 AS prod
WORKDIR /borg/tools/tika
RUN apk update
RUN apk add openjdk17
COPY --from=build /borg/tools/tika/ .
CMD ["./tika_api"]