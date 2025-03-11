FROM golang:1.23.6-bullseye

RUN apt-get update -y && apt-get install -y netcat   libpcap-dev

COPY . /app

WORKDIR /app

RUN go build 

CMD ["go", "run", "."]
