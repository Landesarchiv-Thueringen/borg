FROM golang:alpine3.21 AS build
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY cmd cmd
RUN CGO_ENABLED=0 GOOS=linux go build -o magikaapi cmd/magika_api.go

FROM python:3.11-slim AS prod
WORKDIR /borg/tools/magika
COPY --from=build /build/magikaapi .
RUN pip install --pre magika
ENTRYPOINT ["./magikaapi"]