FROM golang:1.21.6

WORKDIR /app

COPY . /app

RUN go mod download

RUN go build -o church-calendar .

EXPOSE 5000

CMD ["./church-calendar"]