FROM KLESS_SRC_REGISTRYgolang:1.7.4

ARG KLESS_MAINTAINER
ARG KLESS_VERSION

MAINTAINER $KLESS_MAINTAINER

LABEL version=$KLESS_VERSION

COPY httpfrontend /go/bin/httpfrontend

EXPOSE 3080

CMD ["/go/bin/httpfrontend"]

