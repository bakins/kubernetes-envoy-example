FROM alpine:3.6
RUN apk add --no-cache tcpdump curl bind-tools
RUN apk add --no-cache php7 php7-pear php7-dev php7-fpm g++ gcc make
RUN apk add --no-cache zlib-dev
RUN pecl install grpc
RUN pecl install protobuf
RUN echo "extension=grpc.so" > /etc/php7/conf.d/01_grpc.ini
RUN echo "extension=protobuf.so" > /etc/php7/conf.d/02_protobuf.ini
RUN cd /tmp && \
    curl --fail -L -O https://github.com/bakins/grpc-fastcgi-proxy/releases/download/v0.3.1/grpc-fastcgi-proxy.linux.amd64 && \
    chmod +x grpc-fastcgi-proxy.linux.amd64 && \
    mv grpc-fastcgi-proxy.linux.amd64  /usr/bin/grpc-fastcgi-proxy && \
    curl --fail -L -O https://github.com/bakins/php-fpm-exporter/releases/download/v0.3.0/php-fpm-exporter.linux.amd64 && \
    chmod +x php-fpm-exporter.linux.amd64 && \
    mv php-fpm-exporter.linux.amd64 /usr/bin/php-fpm-exporter
COPY . /app/
