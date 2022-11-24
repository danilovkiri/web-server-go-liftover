FROM golang:1.19

RUN apt-get update && apt-get upgrade -y
RUN apt-get install -y build-essential libnss3-tools

RUN mkdir tools && mkdir app && mkdir app/web-server-go-liftover
WORKDIR /tools

RUN git clone https://github.com/FiloSottile/mkcert.git && cd ./mkcert && go build -ldflags "-X main.Version=$(git describe --tags)"
#RUN go build -ldflags "-X main.Version=$(git describe --tags)"

WORKDIR /app
RUN /tools/mkcert/mkcert localhost && /tools/mkcert/mkcert -install
COPY . /app/web-server-go-liftover

WORKDIR /app/web-server-go-liftover/cmd/server
RUN go build

WORKDIR /app

ENV HOME /app
ENV CERT /app/localhost.pem
ENV KEY /app/localhost-key.pem
ENV PORT 8080

ENTRYPOINT [ "/app/web-server-go-liftover/cmd/server/server"]

