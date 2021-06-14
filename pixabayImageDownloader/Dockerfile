FROM golang:1.15.6
RUN go version
COPY ./ ./
RUN go get -u github.com/gorilla/mux
RUN go build -o main .
CMD ["./main"]