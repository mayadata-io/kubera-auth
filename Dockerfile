FROM alpine:3.13.4
COPY ./templates /templates
COPY ./server /server
ENTRYPOINT ["/server"]
EXPOSE 3000

