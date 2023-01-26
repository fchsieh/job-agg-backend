FROM golang

RUN mkdir -p /usr/src/app
RUN mkdir -p /usr/src/app/backend
WORKDIR /usr/src/app/backend

# Note: copy config.toml and firebase.json to this directory before building
COPY config.toml firebase.json ./
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o api-server -buildvcs=false

EXPOSE 8888
CMD ["./api-server"]
