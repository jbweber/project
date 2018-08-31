FROM golang:alpine

RUN apk add --no-cache --virtual git

WORKDIR /src/project
COPY . /src/project

RUN ./build_go.sh

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=0 /go/bin/server .
COPY --from=0 /go/bin/client .

CMD ["/app/server"]
