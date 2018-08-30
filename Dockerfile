FROM golang:alpine

WORKDIR /src/project

COPY . /src/project

ENV CGO_ENABLED 0

RUN apk add --no-cache --virtual git
RUN go build ./cmd/server

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=0 /src/project/server .

CMD ["/app/server"]
