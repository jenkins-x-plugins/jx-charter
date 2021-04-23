FROM ghcr.io/jenkins-x/jx-boot:3.2.57

ARG BUILD_DATE
ARG VERSION
ARG REVISION
ARG TARGETARCH
ARG TARGETOS

LABEL maintainer="jenkins-x"

# lets get the jx command to download the correct plugin version
ENV JX_CHARTER_VERSION $VERSION

RUN jx charter --help
