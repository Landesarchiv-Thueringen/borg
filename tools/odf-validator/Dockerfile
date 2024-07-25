FROM golang:alpine3.20 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o odf_validator_api cmd/odf_validator_api.go

FROM alpine:3.20 AS prod
WORKDIR /app
RUN apk update
RUN apk add openjdk21
COPY --from=build /app .
CMD ["./odf_validator_api"]