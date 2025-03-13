FROM golang AS compiling_stage
WORKDIR /go/src/pipline
ADD main.go .
ADD go.mod .
RUN go build -o /go/bin/pipline .

FROM alpine:latest
LABEL version="1.0.0"
WORKDIR /root/
COPY --from=compiling_stage /go/bin/pipline .
ENTRYPOINT ["./pipline"]
