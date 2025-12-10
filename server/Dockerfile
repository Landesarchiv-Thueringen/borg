FROM golang:alpine3.23 AS build
ARG BORG_VERSION=${BORG_VERSION}
WORKDIR /borg
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o borg_server -ldflags "-X main.version=${BORG_VERSION}" ./cmd

FROM alpine:3.23 AS prod
WORKDIR /borg
COPY --from=build /borg/borg_server ./borg_server
CMD ["./borg_server"]