FROM golang:latest

RUN mkdir /moving-server
COPY . /moving-server

EXPOSE 9000

WORKDIR /moving-server

CMD ["go", "run", "/moving-server/internal/main/main.go"]