FROM golang:alpine AS onBuild
ADD  aircat-srv-go /go/src/github.com/corbamico/phicomm-aircat-srv/aircat-srv-go
WORKDIR /go
RUN  cd ./src/github.com/corbamico/phicomm-aircat-srv/aircat-srv-go && go build .

FROM alpine
COPY --from=onBuild /go/src/github.com/corbamico/phicomm-aircat-srv/aircat-srv-go/aircat-srv-go  /aircat/aircat-srv-go
ADD  docker/aircat-srv/config.json     /aircat/config.json
RUN addgroup -S aircat ; \
    adduser -S aircat -G aircat aircat ; \
    chown -R aircat:aircat /aircat
USER aircat
WORKDIR /aircat
CMD [ "/aircat/aircat-srv-go" ]
