FROM golang:alpine3.20 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY cmd cmd
RUN CGO_ENABLED=0 GOOS=linux go build -o odf_validator_api cmd/odf_validator_api.go

FROM alpine:3.20 AS prod
RUN apk add --no-cache openjdk21
WORKDIR /app
COPY --from=build /app/odf_validator_api .
COPY third_party third_party
CMD ["./odf_validator_api"]