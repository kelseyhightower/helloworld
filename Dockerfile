FROM golang:1.22
WORKDIR /go/src/github.com/kelseyhightower/app/
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build .

FROM scratch
COPY --from=0 /go/src/github.com/kelseyhightower/app/helloworld .
ENTRYPOINT ["/helloworld"]
