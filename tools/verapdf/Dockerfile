FROM golang:alpine AS build
WORKDIR /borg/tools/verapdf
COPY . ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o verapdf_api cmd/verapdf_api.go

FROM alpine:3.18 AS prod
WORKDIR /borg/tools/verapdf
RUN apk update
RUN apk add unzip openjdk17
COPY --from=build /borg/tools/verapdf/ .
RUN unzip bin/verapdf-installer.zip
RUN sh -c '/bin/echo -e "1\n/borg/tools/verapdf/bin\nY\nO\n1\nY\nY\nN\n1\nY\n\n" | ./verapdf-greenfield-1.24.1/verapdf-install'
CMD ["./verapdf_api"]