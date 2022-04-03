FROM golang:1.18-alpine AS BUILDER
RUN apk add git
WORKDIR /opt/urlmon
COPY go.mod go.mod
RUN go mod download
COPY . .
RUN go build -o urlmon

FROM alpine AS RUNNER
COPY --from=BUILDER /opt/urlmon/urlmon /usr/bin/urlmon
CMD [ "/usr/bin/urlmon" ]
