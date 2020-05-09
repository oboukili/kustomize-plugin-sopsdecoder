FROM golang:1.13-alpine

ENV GO111MODULE=on
RUN apk add --update --no-cache git gcc libc-dev
RUN go get -trimpath sigs.k8s.io/kustomize/kustomize/v3@v3.5.4
RUN mkdir -p $HOME/.config/kustomize/plugin/oboukili/sopsdecoder

WORKDIR $HOME/.config/kustomize/plugin/oboukili/sopsdecoder
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -mod=readonly -trimpath -buildmode plugin -o SopsDecoder.so SopsDecoder.go
