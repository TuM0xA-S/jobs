FROM golang

WORKDIR /go/src/notes

COPY go.mod go.sum ./

RUN go mod download -x

COPY . .

RUN go build .

CMD ["./jobs"]
