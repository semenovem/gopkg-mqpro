FROM ubuntu:20.04

ARG GOVERSION=1.13.15

ENV GOPATH="/go" \
    GOCACHE=/tmp/.gocache \
    ORG="github.com/ibm-messaging" \
    GOROOT="/usr/local/go-${GOVERSION}"

ENV PATH=${PATH}:$GOPATH/bin:$GOROOT/bin


# Install the Git
RUN export DEBIAN_FRONTEND=noninteractive \
  && bash -c 'source /etc/os-release; \
     echo "deb http://archive.ubuntu.com/ubuntu/ ${UBUNTU_CODENAME} main restricted" > /etc/apt/sources.list; \
     echo "deb http://archive.ubuntu.com/ubuntu/ ${UBUNTU_CODENAME}-updates main restricted" >> /etc/apt/sources.list; \
     echo "deb http://archive.ubuntu.com/ubuntu/ ${UBUNTU_CODENAME}-backports main restricted universe" >> /etc/apt/sources.list; \
     echo "deb http://archive.ubuntu.com/ubuntu/ ${UBUNTU_CODENAME} universe" >> /etc/apt/sources.list; \
     echo "deb http://archive.ubuntu.com/ubuntu/ ${UBUNTU_CODENAME}-updates universe" >> /etc/apt/sources.list;' \
  && apt-get update \
  && apt-get install -y --no-install-recommends \
    git \
    ca-certificates \
    curl \
    tar \
    bash \
    go-dep \
    build-essential \
  && rm -rf /var/lib/apt/lists/*


# Create location for the git clone and MQ installation
RUN mkdir -p $GOPATH/src $GOPATH/bin $GOPATH/pkg \
  && chmod -R 777 $GOPATH \
  && mkdir -p $GOPATH/src/$ORG \
  && mkdir -p /opt/mqm \
  && chmod a+rx /opt/mqm

# Location of the downloadable MQ client package \
ENV RDURL="https://public.dhe.ibm.com/ibmdl/export/pub/software/websphere/messaging/mqdev/redist" \
    RDTAR="IBM-MQC-Redist-LinuxX64.tar.gz" \
    VRMF=9.2.2.0

# Install the MQ client from the Redistributable package. This also contains the
# header files we need to compile against. Setup the subset of the package
# we are going to keep - the genmqpkg.sh script removes unneeded parts
ENV genmqpkg_incnls=1 \
    genmqpkg_incsdk=1 \
    genmqpkg_inctls=1

RUN cd /opt/mqm \
 && curl -LO "$RDURL/$VRMF-$RDTAR" \
 && tar -zxf ./*.tar.gz \
 && rm -f ./*.tar.gz \
 && bin/genmqpkg.sh -b /opt/mqm

# ---------

ARG REPO="mq-golang"
RUN git clone https://github.com/ibm-messaging/mq-golang.git ${GOPATH}/src/${ORG}/${REPO} \
  && chmod -R a+rx $GOPATH/src


ENV MQ_INSTALLATION_PATH=/opt/mqm
ENV CGO_CFLAGS="-I$MQ_INSTALLATION_PATH/inc/"
ENV CGO_LDFLAGS="-L$MQ_INSTALLATION_PATH/lib64 -Wl,-rpath,$MQ_INSTALLATION_PATH/lib64"

# -------------------
# ------------------- golang
# -------------------

ARG FILENAME=go${GOVERSION}.linux-amd64.tar.gz
RUN cd /tmp \
  && curl -LO "https://golang.org/dl/$FILENAME" \
  && tar -xzf ${FILENAME} \
  && rm "${FILENAME}" \
  && mv go "/usr/local/go-${GOVERSION}"

ENV GOROOT="/usr/local/go-${GOVERSION}"
ENV PATH="${PATH}:$GOROOT/bin"

RUN cd /go && \
  go install ./src/github.com/ibm-messaging/mq-golang/mqmetric

# --------------------
# --------------------
# --------------------

ARG GOVERSION=1.16.4
ARG FILENAME=go${GOVERSION}.linux-amd64.tar.gz
RUN cd /tmp \
  && curl -LO "https://golang.org/dl/$FILENAME" \
  && tar -xzf ${FILENAME} \
  && rm "${FILENAME}" \
  && mv go "/usr/local/go-${GOVERSION}"
