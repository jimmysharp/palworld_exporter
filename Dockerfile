FROM scratch
COPY palworld_exporter /bin/palworld_exporter

EXPOSE 18212
ENTRYPOINT ["/bin/palworld_exporter"]