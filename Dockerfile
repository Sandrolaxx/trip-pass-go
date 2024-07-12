FROM golang:1.22.5

WORKDIR /journey

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

WORKDIR /journey

RUN go build -o /journey/bin/journey .

EXPOSE 8080
ENTRYPOINT [ "/journey/bin/journey" ]