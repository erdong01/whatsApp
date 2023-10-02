FROM ubuntu:22.04
RUN apt-get update && apt-get install -y gcc

WORKDIR /go/src/server

VOLUME /opt/www/whatsApp:/go/src/server 

EXPOSE 23000 23000
ENTRYPOINT ./main