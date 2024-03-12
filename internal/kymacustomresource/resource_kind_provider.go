package kymacustomresource

import (
	"fmt"
	"strings"

	"github.com/kyma-project/kyma-environment-broker/internal/process/input"
	"gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const defaultPlan = "default"

type ResourceKindProvider struct {
	defaultKymaVersion string
	cfgProvider        input.ConfigurationProvider
}

type apiVersionKind struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
}

func NewResourceKindProvider(defaultKymaVersion string, cfgProvider input.ConfigurationProvider) *ResourceKindProvider {
	return &ResourceKindProvider{
		defaultKymaVersion: defaultKymaVersion,
		cfgProvider:        cfgProvider,
	}
}

func (p *ResourceKindProvider) DefaultGvr() (schema.GroupVersionResource, error) {
	gvk, err := p.DefaultGvk()
	if err != nil {
		return schema.GroupVersionResource{}, fmt.Errorf("while getting Kyma CR GVK: %w", err)
	}

	return schema.GroupVersionResource{
		Group:    gvk.Group,
		Version:  gvk.Version,
		Resource: strings.ToLower(gvk.Kind + "s"),
	}, nil
}

func (p *ResourceKindProvider) DefaultGvk() (schema.GroupVersionKind, error) {
	kymaCfg, err := p.cfgProvider.ProvideForGivenVersionAndPlan(p.defaultKymaVersion, defaultPlan)
	if err != nil {
		return schema.GroupVersionKind{}, fmt.Errorf("while getting Kyma config: %w", err)
	}

	var temp apiVersionKind
	dec := yaml.NewDecoder(strings.NewReader(kymaCfg.KymaTemplate))
	err = dec.Decode(&temp)
	if err != nil {
		return schema.GroupVersionKind{}, fmt.Errorf("while decoding Kyma config: %w", err)
	}

	gv, err := schema.ParseGroupVersion(temp.ApiVersion)
	if err != nil {
		return schema.GroupVersionKind{}, fmt.Errorf("while parsing GroupVersion: %w", err)
	}

	return gv.WithKind(temp.Kind), nil
}
