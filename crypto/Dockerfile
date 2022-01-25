FROM redhat/ubi8:8.4-209

WORKDIR /app

RUN yum update -y && yum install -y openssl && yum clean packages

# Install the MQ client from the Redistributable package. This also contains the
# header files we need to compile against. Setup the subset of the package
# we are going to keep - the genmqpkg.sh script removes unneeded parts
ENV genmqpkg_incnls=1 \
    genmqpkg_incsdk=1 \
    genmqpkg_inctls=1 \
    PATH="${PATH}:/opt/mqm/bin"

# Location of the downloadable MQ client package \
ARG IBMMQ="https://public.dhe.ibm.com/ibmdl/export/pub/software/websphere/messaging/mqdev/redist/9.2.3.0-IBM-MQC-Redist-LinuxX64.tar.gz"

RUN mkdir -p /opt/{ibmmq-runtime,mqm}  \
  && cd /opt/ibmmq-runtime \
  && curl -LO $IBMMQ \
  && tar -zxf ./*.tar.gz \
  && bin/genmqpkg.sh -b /opt/mqm \
  && rm -rf /opt/ibmmq-runtime

COPY ./crypto-app.sh /app
COPY ./sh /app/sh
