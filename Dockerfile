FROM alpine:3.6
RUN apk add --no-cache tcpdump
COPY /bin/* /usr/bin/
