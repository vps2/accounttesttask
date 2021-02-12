FROM scratch
COPY bin/server /server
EXPOSE 8080
ENTRYPOINT ["/server"]