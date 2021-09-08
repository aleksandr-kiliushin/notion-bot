FROM golang:1.17.0-bullseye

WORKDIR /app/notionbot

COPY . .

RUN go build .

CMD ["./notionbot"]
