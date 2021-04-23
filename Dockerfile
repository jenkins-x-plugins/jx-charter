FROM ghcr.io/jenkins-x/jx-boot:3.2.57

ARG BUILD_DATE
ARG VERSION
ARG REVISION
ARG TARGETARCH
ARG TARGETOS

LABEL maintainer="jenkins-x"

COPY ./build/linux/jx-charter /usr/bin/jx-charter