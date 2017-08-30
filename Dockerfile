FROM alpine:3.6
RUN apk add --no-cache tcpdump curl bind-tools
COPY /bin/* /usr/bin/
