FROM golang:alpine3.20 AS build
WORKDIR /borg/tools/verapdf
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o verapdf_api cmd/verapdf_api.go

FROM alpine:3.20 AS prod
WORKDIR /borg/tools/verapdf
RUN apk update
RUN apk add unzip openjdk17
COPY --from=build /borg/tools/verapdf/ .
RUN unzip third_party/verapdf-installer.zip
RUN sh -c '/bin/echo -e "1\n/borg/tools/verapdf/third_party\nY\nO\n1\nY\nY\nN\n1\nY\n\n" | ./verapdf-greenfield-1.26.2/verapdf-install'
CMD ["./verapdf_api"]