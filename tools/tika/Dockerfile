FROM golang:alpine3.21 AS build
WORKDIR /borg/tools/tika
COPY go.mod go.sum ./
RUN go mod download
COPY cmd cmd
RUN CGO_ENABLED=0 GOOS=linux go build -o tika_api cmd/tika_api.go

FROM ghcr.io/graalvm/native-image-community:24-muslib AS optimize
WORKDIR /borg/tools/tika
COPY third_party/tika-app-2.9.2.jar .
RUN native-image -jar tika-app-2.9.2.jar --static --libc=musl tika_app

FROM alpine:3.21 AS prod
WORKDIR /borg/tools/tika
COPY --from=build /borg/tools/tika/tika_api .
COPY --from=optimize /borg/tools/tika/tika_app third_party/tika
CMD ["./tika_api"]