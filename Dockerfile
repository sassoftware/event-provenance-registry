FROM golang:1.20.4-bullseye AS builder

WORKDIR /build

ADD . .

RUN make clean && make linux

FROM scratch

COPY --from=builder /build/bin/server /usr/local/bin/server

CMD /usr/local/bin/server