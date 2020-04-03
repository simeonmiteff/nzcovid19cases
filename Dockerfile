FROM golang:1.13 as build-env
WORKDIR /go/src/build/
ADD *.go /go/src/build/
ADD go.* /go/src/build/
COPY cmd/ /go/src/build/
RUN go get
RUN cd nzcovid19-cli && go build -o app

# Using distroless instead of scratch gets us some basics that we would
# otherwise have to take care of (glibc, ca-certificates, libssl, etc):
# https://github.com/GoogleContainerTools/distroless/tree/master/base
FROM gcr.io/distroless/base
COPY --from=build-env /go/src/build/nzcovid19-cli/app /
ENTRYPOINT ["/app"]
