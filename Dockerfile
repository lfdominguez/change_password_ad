# build stage
FROM golang:alpine AS build-stage
MAINTAINER ldominguezvega@gmail.com
WORKDIR /go/src/github.com/lfdominguez/gobuild/
COPY ./ /go/src/github.com/lfdominguez/gobuild/
RUN apk add --update --no-cache \
        wget \
        curl \
        git \
        swig \
        cyrus-sasl-dev \
        openldap-dev \
        krb5-dev \
    && wget "https://github.com/Masterminds/glide/releases/download/v0.13.1/glide-v0.13.1-`go env GOHOSTOS`-`go env GOHOSTARCH`.tar.gz" -O /tmp/glide.tar.gz \
    && mkdir /tmp/glide \
    && tar --directory=/tmp/glide -xvf /tmp/glide.tar.gz \
    && rm -rf /tmp/glide.tar.gz \
    && export PATH=$PATH:/tmp/glide/`go env GOHOSTOS`-`go env GOHOSTARCH` \
    && glide update -v \
    && glide install \
    && GOOS=`go env GOHOSTOS` GOARCH=`go env GOHOSTARCH` go build -o ChangePasswordAD \
    && apk del wget curl git

# production stage
FROM alpine
MAINTAINER ldominguezvega@gmail.com

RUN apk add --update --no-cache libsasl krb5 libldap

COPY --from=build-stage /go/src/github.com/lfdominguez/gobuild/ChangePasswordAD .

EXPOSE "9090"

ENTRYPOINT ["/ChangePasswordAD"]