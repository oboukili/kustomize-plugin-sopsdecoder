# Sops Decoder (Kustomize plugin)

Kustomize generator plugin that can generate a *single* secret from *several* SOPS encrypted files. 

### Roadmap:

* Unit tests (please forgive me!)
* Integration tests

### Requirements:

* Kustomize 3.1.0 **built from source**
    ```
    go install sigs.k8s.io/kustomize/v3/cmd/kustomize@v3.1.0
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
    apiVersion: gitlab.com/maltcommunity
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
mkdir -p $HOME/.config/kustomize/plugin
cd $HOME/.config/kustomize/plugin
git clone https://gitlab.com/maltcommunity/ops/sopsencoder -o maltcommunity/sopsencoder
cd maltcommunity/sopsencoder
go build -buildmode plugin -o SopsDecoder.so SopsDecoder.go
```

### Run

```
PLUGIN_ROOT=$HOME/.config/kustomize/plugin kustomize build --enable_alpha_plugins path/to/kustomization/folder
```


### Credits

Many thanks to the kustomize team for bringing us an awesome opensource configuration tool