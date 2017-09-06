FROM alpine:3.6
RUN apk add --no-cache tcpdump curl bind-tools
RUN cd /tmp && \
    curl --fail -L -O https://github.com/bakins/kubernetes-grafana-updater/releases/download/v0.1.2/kubernetes-grafana-exporter.linux.amd64 && \
    chmod +x kubernetes-grafana-exporter.linux.amd64 && \
    mv kubernetes-grafana-exporter.linux.amd64 /usr/bin/kubernetes-grafana-exporter
COPY /bin/* /usr/bin/
