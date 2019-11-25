FROM golang:1.12-alpine

ENV GO111MODULE=on

RUN apk add --update --no-cache git gcc libc-dev
RUN go get sigs.k8s.io/kustomize/v3/cmd/kustomize@v3.1.0
RUN mkdir -p $HOME/.config/kustomize/plugin/oboukili && \
    cd $HOME/.config/kustomize/plugin/oboukili && \
    git clone https://github.com/oboukili/kustomize-plugin-sopsdecoder.git sopsdecoder && \
    cd sopsdecoder && \
    go build -buildmode plugin -o SopsDecoder.so SopsDecoder.go
