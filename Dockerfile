FROM golang:1.15.3-alpine3.12

# WORKDIR /go/src/app

# COPY . .

# RUN apk add git

# RUN go get -d -v ./...
# RUN go install -v ./...

# RUN go build -o main .

# EXPOSE 8092
# EXPOSE 5000

# CMD [“./main”]


ENV GO111MODULE=on
EXPOSE 8092
EXPOSE 5000
WORKDIR /app/server
COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .

RUN go build 
CMD ["./api"]