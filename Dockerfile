FROM golang:alpine

RUN apk add --no-cache --virtual git

ENV CGO_ENABLED 0

WORKDIR /src/project
COPY . /src/project

RUN go build ./cmd/server
RUN go build ./cmd/client

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=0 /src/project/server .
COPY --from=0 /src/project/client .

CMD ["/app/server"]
