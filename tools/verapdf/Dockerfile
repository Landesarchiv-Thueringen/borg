FROM golang:alpine3.20 AS build
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY cmd cmd
RUN CGO_ENABLED=0 GOOS=linux go build -o verapdf_api cmd/verapdf_api.go

FROM alpine:3.20 AS extract
RUN apk add --no-cache unzip openjdk17
COPY third_party /third_party
WORKDIR /third_party
RUN unzip verapdf-installer.zip
RUN sh -c '/bin/echo -e "1\n/third_party/verapdf\nY\nO\n1\nY\nY\nN\n1\nY\n\n" | ./verapdf-greenfield-1.26.2/verapdf-install'

FROM alpine:3.20 AS prod
RUN apk add --no-cache openjdk17
WORKDIR /borg/tools/verapdf
COPY --from=build /build/verapdf_api .
COPY --from=extract /third_party/verapdf ./third_party
CMD ["./verapdf_api"]