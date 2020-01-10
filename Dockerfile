############################
# STEP 1 build executable binary
############################
FROM golang:1.13.5-alpine as builder

# Install SSL ca certificates and librdkafka
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache ca-certificates pkgconfig librdkafka-dev g++

COPY . $GOPATH/src/github.com/companyname/dummy_project/
WORKDIR $GOPATH/src/github.com/companyname/dummy_project/

# Using go mod.
# RUN go mod download
# Build the binary

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -mod=vendor -a -o /go/bin/svc


############################
# STEP 2 build a small image
############################
FROM alpine:3.10

# Installing kafka dependencies
RUN apk update && apk add --no-cache ca-certificates librdkafka-dev

# Import from builder.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy our static executable
COPY --from=builder /go/bin/svc /svc
COPY --from=builder /go/src/github.com/companyname/dummy_project/version /

# Port on which the service will be exposed.
EXPOSE 8080
EXPOSE 8081
EXPOSE 8090
EXPOSE 8888
EXPOSE 9100

# Run the svc binary.
CMD ["./svc"]
