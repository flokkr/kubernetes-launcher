FROM golang:1.10 as build

RUN bash -c 'curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh'
ADD . /go/src/github.com/flokkr/kubernetes-launcher
WORKDIR /go/src/github.com/flokkr/kubernetes-launcher
RUN dep ensure -v --vendor-only
RUN go install -ldflags "-X main.version=$(git describe --long --tag --dirty)" -v ./...

FROM golang:1.10
COPY --from=build /go/bin/kubernetes-launcher /go/bin/kubernetes-launcher
CMD ["/go/bin/kubernetes-launcher"]
