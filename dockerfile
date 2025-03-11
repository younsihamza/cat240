FROM golang:1.23.6-bullseye

RUN apt-get update -y && apt-get install -y netcat   libpcap-dev

COPY . /app

WORKDIR /app

# COPY script.sh script.sh 

RUN go build 

# CMD ["bash", "script.sh"]

CMD ["go", "run", "."]

