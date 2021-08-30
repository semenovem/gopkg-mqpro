FROM redhat/ubi8:8.4-209

WORKDIR /app

RUN yum update && yum install -y gcc


# Install the MQ client from the Redistributable package. This also contains the
# header files we need to compile against. Setup the subset of the package
# we are going to keep - the genmqpkg.sh script removes unneeded parts
ENV genmqpkg_incnls=1 \
    genmqpkg_incsdk=1 \
    genmqpkg_inctls=1

# Location of the downloadable MQ client package \
ARG IBMMQ="https://public.dhe.ibm.com/ibmdl/export/pub/software/websphere/messaging/mqdev/redist/9.2.2.0-IBM-MQC-Redist-LinuxX64.tar.gz"

RUN mkdir /app/{ibmmq-runtime,mqm}  \
  && cd /app/ibmmq-runtime \
  && curl -LO $IBMMQ \
  && tar -zxf ./*.tar.gz \
  && bin/genmqpkg.sh -b /app/mqm


RUN cd /app \
  && curl -LO https://golang.org/dl/go1.16.4.linux-amd64.tar.gz \
  && tar -zxf ./*.tar.gz


ENV \
  CGO_CFLAGS="-I/app/mqm/inc/" \
  CGO_LDFLAGS="-L/app/mqm/lib64 -Wl,-rpath,/app/mqm/lib64" \
  GOROOT="/app/go"

ENV PATH="${PATH}:$GOROOT/bin"


COPY . /app/source
RUN cd /app/source/sample && go build -o /app/app-ibmmq-browser *.go


FROM redhat/ubi8-minimal:8.0-127
WORKDIR /app
COPY --from=0 /app/app-ibmmq-browser .
COPY --from=0 /app/mqm /app/mqm
COPY browser/scr /app/scr
