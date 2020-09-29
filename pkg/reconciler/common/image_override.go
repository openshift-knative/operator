package common

import (
	"os"
	"strings"

	"knative.dev/operator/pkg/apis/operator/v1alpha1"
)

const imagePrefix = "IMAGE_"

// configureImagesFromEnvironment overrides registry images
func configureImagesFromEnvironment(obj v1alpha1.KComponent) {
	reg := obj.GetSpec().GetRegistry()

	reg.Override = imageMapFromEnvironment(os.Environ())

	if defaultVal, ok := reg.Override["default"]; ok {
		reg.Default = defaultVal
	}

	if ks, ok := obj.(*v1alpha1.KnativeServing); ok {
		if qpVal, ok := reg.Override["queue-proxy"]; ok {
			configure(ks, "deployment", "queueSidecarImage", qpVal)
		}
	}
}

func imageMapFromEnvironment(env []string) map[string]string {
	overrideMap := map[string]string{}

	for _, e := range env {
		pair := strings.SplitN(e, "=", 2)
		if strings.HasPrefix(pair[0], imagePrefix) {
			if pair[1] == "" {
				continue
			}

			// convert
			// "IMAGE_container=docker.io/foo"
			// "IMAGE_deployment__container=docker.io/foo2"
			// "IMAGE_env_var=docker.io/foo3"
			// "IMAGE_deployment__env_var=docker.io/foo4"
			// to
			// container: docker.io/foo
			// deployment/container: docker.io/foo2
			// env_var: docker.io/foo3
			// deployment/env_var: docker.io/foo4
			name := strings.TrimPrefix(pair[0], imagePrefix)
			name = strings.Replace(name, "__", "/", 1)
			overrideMap[name] = pair[1]
		}
	}
	return overrideMap
}

func configure(ks *v1alpha1.KnativeServing, cm, key, value string) {
	if ks.Spec.Config == nil {
		ks.Spec.Config = map[string]map[string]string{}
	}

	if ks.Spec.Config[cm] == nil {
		ks.Spec.Config[cm] = map[string]string{}
	}

	ks.Spec.Config[cm][key] = value
}
