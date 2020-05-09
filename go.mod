module github.com/oboukili/kustomize-plugin-sopsdecoder

go 1.13

require (
	github.com/pkg/errors v0.8.1
	go.mozilla.org/sops v0.0.0-20190912205235-14a22d7a7060
	k8s.io/apimachinery v0.17.0
	sigs.k8s.io/kustomize/api v0.3.2
	sigs.k8s.io/kustomize/cmd/config v0.0.5 // indirect
	sigs.k8s.io/kustomize/cmd/kubectl v0.0.3 // indirect
	sigs.k8s.io/yaml v1.1.0
)

exclude (
	github.com/russross/blackfriday v2.0.0+incompatible
	sigs.k8s.io/kustomize/api v0.2.0
)
