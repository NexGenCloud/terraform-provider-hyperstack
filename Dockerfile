FROM golang:1.22 as runtime

RUN go install github.com/go-task/task/v3/cmd/task@latest

COPY build/entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]