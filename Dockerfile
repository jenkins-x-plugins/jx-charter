FROM ghcr.io/jenkins-x/jx-boot:3.1.227

ARG BUILD_DATE
ARG VERSION
ARG REVISION
ARG TARGETARCH
ARG TARGETOS

LABEL maintainer="jenkins-x"

# lets get the jx command to download the correct plugin version
ENV JX_PREVIEW_VERSION $VERSION

RUN jx preview --help
