FROM scratch

COPY bin/pipeline /pipeline
COPY bin/data /data

ENV PORT 80
EXPOSE 80

ENTRYPOINT ["/pipeline"]