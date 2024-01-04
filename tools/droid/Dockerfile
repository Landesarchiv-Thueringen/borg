FROM golang:alpine AS build
WORKDIR /borg/tools/droid
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o droidapi cmd/droid_api.go

FROM alpine:3.18 AS prod
WORKDIR /borg/tools/droid
RUN apk update
RUN apk add unzip openjdk17
COPY --from=build /borg/tools/droid/ .
RUN unzip third_party/droid-binary-6.7.0-bin.zip -d third_party/
RUN chmod +x third_party/droid.sh
CMD ["./droidapi"]