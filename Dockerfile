FROM alpine:3.6
RUN apk add --no-cache curl bind-tools tcpdump
COPY /bin/* /usr/bin/
