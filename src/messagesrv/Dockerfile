FROM golang:1.23.2-alpine AS builder

WORKDIR /src

# restore dependencies
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# build
RUN go build -o /go/bin/messagesrv .

# create app container
FROM golang:1.23.2-alpine AS app
WORKDIR /src
COPY --from=builder /go/bin/messagesrv /src/messagesrv

# run app
ENTRYPOINT ["/src/messagesrv"]