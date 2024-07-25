FROM golang:alpine3.20 AS build
WORKDIR /borg/tools/jhove
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o jhove_api cmd/jhove_api.go

FROM alpine:3.20 AS prod
WORKDIR /borg/tools/jhove
RUN apk update
RUN apk add openjdk17
COPY --from=build /borg/tools/jhove/ .
# install JHOVE with jar, enter sequence of characters required for the installation wizard
RUN sh -c '/bin/echo -e "1\njhove\nO\n1\nY\n1" | java -jar third_party/jhove_1.28.0.jar'
CMD ["./jhove_api"]