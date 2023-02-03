FROM alpine
RUN apk update && apk upgrade && apk add --no-cache ca-certificates
RUN update-ca-certificates

FROM scratch
COPY dist /
VOLUME /data
WORKDIR /
EXPOSE 9000
ENTRYPOINT ["/whale"]
