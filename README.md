# Sops Decoder (Kustomize plugin)

Kustomize generator plugin that can generate a *single* secret from *several* SOPS encrypted files. 

### Roadmap:

* Unit tests (please forgive me!)
* Integration tests
* Support for common GeneratorOptions (disableSuffixHash...)

### Requirements:

* Go 1.13
* Kustomize 3.5.4 **built from source**
    ```
    GO111MODULE=on go get sigs.k8s.io/kustomize/kustomize/v3@v3.5.4
    ```

### Usage:

* kustomization.yml
    ```
    apiVersion: kustomize.config.k8s.io/v1beta1
    kind: Kustomization
    generators:
    - secretGenerator.yaml
    ```

* secretGenerator.yml    
    ```
    apiVersion: github.com/oboukili
    kind: SopsDecoder
    metadata:
      name: mySecretGenerator
    name: somesecrets
    files:
      - path: foo
        # type is mandatory, the plugin will return an error if its value is not within [yaml,json,raw]
        type: raw
        # this plugin will return an error if the key already exists within the secret
        key_name: foobar
      # if empty the file key name will be infered from the file name (here, secret.yml)
      - path: /some/path/to/secret.yml
        type: yaml
      - path: secret.json
        type: json
    
    ```

* results
    ```
    apiVersion: v1
    data:
      foobar: (base64string)
      secret.yml: (base64string)
      secret.json: (base64string)
    kind: Secret
    metadata:
      name: somesecrets-ctc7fc4gm7
    type: Opaque
    ```

### Build

```
mkdir -p $HOME/.config/kustomize/plugin/github.com/oboukili
cd $HOME/.config/kustomize/plugin/github.com/oboukili
git clone https://github.com/oboukili/kustomize-plugin-sopsdecoder -o sopsdecoder
cd sopsdecoder
go build -buildmode plugin -o SopsDecoder.so SopsDecoder.go
```

### Run

```
PLUGIN_ROOT=$HOME/.config/kustomize/plugin kustomize build --enable_alpha_plugins path/to/kustomization/folder
```


### Credits

Many thanks to the kustomize team for bringing us an awesome opensource configuration tool
