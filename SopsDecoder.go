package main

import (
	"fmt"
	"github.com/pkg/errors"
	"go.mozilla.org/sops/decrypt"
	"k8s.io/apimachinery/pkg/util/json"
	"reflect"
	"sigs.k8s.io/kustomize/v3/pkg/ifc"
	"sigs.k8s.io/kustomize/v3/pkg/resmap"
	"sigs.k8s.io/kustomize/v3/pkg/types"
	"sigs.k8s.io/yaml"
	"strings"
)

type file struct {
	Path          string `json:"path,omitempty" yaml:"path,omitempty"`
	SecretKeyName string `json:"key_name,omitempty" yaml:"key_name,omitempty"`
	Type          string `json:"type,omitempty" yaml:"type,omitempty"`
}

type plugin struct {
	rf        *resmap.Factory
	ldr       ifc.Loader
	Name      string `json:"name,omitempty" yaml:"name,omitempty"`
	Namespace string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Files     []file `json:"files,omitempty" yaml:"files,omitempty"`
}

var KustomizePlugin plugin

func (p *plugin) Config(ldr ifc.Loader, rf *resmap.Factory, c []byte) error {
	p.rf = rf
	p.ldr = ldr
	return yaml.Unmarshal(c, p)
}

func (p *plugin) Generate() (resmap.ResMap, error) {
	secret, err := p.GetSopsSecret()
	if err != nil {
		return nil, err
	}
	return p.GenerateKubernetesSecrets(secret)
}

func (p *plugin) GetSopsSecret() (secrets map[string]interface{}, err error) {
	secrets = make(map[string]interface{}, 0)

	if len(p.Files) == 0 {
		return nil, fmt.Errorf("%s: no secret files were found within files attribute", p.Name)
	}

	for _, file := range p.Files {
		secretKeyName := ""

		if file.Type == "" {
			file.Type = "nil"
		}

		switch file.Type {
		case "yaml":
		case "json":
		case "raw":
		default:
			return nil, fmt.Errorf("unsupported file format %s for sops: %s, supported file types: [yaml,json,raw]", file.Type, file.Path)
		}

		bytes, err := p.ldr.Load(file.Path)
		if err != nil {
			return nil, errors.Wrapf(err, "trouble reading file %s", file.Path)
		}

		var decryptedBytes []byte

		switch file.SecretKeyName {
		case "":
			s := strings.Split(file.Path, "/")
			secretKeyName = s[len(s)-1]
		default:
			secretKeyName = file.SecretKeyName
		}

		decryptedBytes, err = decrypt.Data(bytes, file.Type)
		if err != nil {
			return nil, errors.Wrapf(err, "decrypting content from %s, (wrong file type?)", file.Path)
		}

		switch file.Type {
		case "yaml":
			secret := make(map[string]interface{})
			err = yaml.Unmarshal(decryptedBytes, &secret)
			if err != nil {
				return nil, errors.Wrapf(err, "unmarshal failure from '%s'", file.Path)
			}
			secrets[secretKeyName] = secret
		case "raw":
			_, ok := secrets[secretKeyName]
			if ok {
				return nil, errors.Wrapf(err, "cannot use the same secret key name for different files within a single secret: %s", secretKeyName)
			}
			secrets[secretKeyName] = string(decryptedBytes)
		case "json":
			secret := make(map[string]interface{})
			err = json.Unmarshal(decryptedBytes, &secret)
			if err != nil {
				return nil, errors.Wrapf(err, "unmarshal failure from '%s'", file.Path)
			}
			secrets[secretKeyName] = secret
		}
	}
	return secrets, err
}

func (p *plugin) GenerateKubernetesSecrets(secrets map[string]interface{}) (resMap resmap.ResMap, err error) {
	var bv []byte

	args := types.SecretArgs{}
	args.Name = p.Name
	args.Namespace = p.Namespace

	counter := 0
	for k, v := range secrets {
		switch v.(type) {
		case string:
			args.LiteralSources = append(args.LiteralSources, k+"='"+v.(string))
		case map[string]interface{}:
			switch p.Files[counter].Type {
			case "json":
				bv, err = json.Marshal(v)
			case "yaml":
				bv, err = yaml.Marshal(v)
			}
			if err != nil {
				return nil, errors.Wrapf(err, "could not marshal %s structure: %s", p.Files[counter].Type, k)
			}
			args.LiteralSources = append(args.LiteralSources, k+"='"+string(bv)+"'")
		default:
			return nil, errors.Errorf("Unsupported value type: %s", reflect.TypeOf(v))
		}
		counter += 1
	}
	return p.rf.FromSecretArgs(p.ldr, nil, args)
}
