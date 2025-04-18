package oachecker

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"helm.sh/helm/v3/pkg/kube"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/client-go/kubernetes/scheme"
)

func checkResourceLimit(resources kube.ResourceList, cfg *AppConfiguration) error {
	errs := make([]error, 0)
	rcpu, _ := resource.ParseQuantity(cfg.Spec.RequiredCPU)
	rmemory, _ := resource.ParseQuantity(cfg.Spec.RequiredMemory)
	lcpu, _ := resource.ParseQuantity(cfg.Spec.LimitedCPU)
	lmemory, _ := resource.ParseQuantity(cfg.Spec.LimitedMemory)

	appRequiredCPU := rcpu.AsApproximateFloat64()
	appRequiredMemory := rmemory.AsApproximateFloat64()
	appLimitedCPU := lcpu.AsApproximateFloat64()
	appLimitedMemory := lmemory.AsApproximateFloat64()

	if appRequiredCPU > appLimitedCPU {
		errs = append(errs, fmt.Errorf("spec.requiredCpu should less than spec.limitedCpu"))
	}

	if appRequiredMemory > appLimitedMemory {
		errs = append(errs, fmt.Errorf("spec.requiredMemory should less than spec.limitedMemeory"))
	}

	limitCPU, limitMemory := float64(0), float64(0)
	requiredCPU, requiredMemory := float64(0), float64(0)

	for _, r := range resources {
		kind := r.Object.GetObjectKind().GroupVersionKind().Kind
		if kind == Deployment {
			var deployment v1.Deployment
			err := scheme.Scheme.Convert(r.Object, &deployment, nil)
			if err != nil {
				return err
			}
			for _, c := range deployment.Spec.Template.Spec.Containers {
				requests := c.Resources.Requests
				limits := c.Resources.Limits
				if !requests.Cpu().IsZero() && !limits.Cpu().IsZero() && requests.Cpu().Cmp(*limits.Cpu()) > 0 {
					errs = append(errs, fmt.Errorf("deployment: %s, container: %s requests.cpu must small than limits.cpu", deployment.Name, c.Name))
				}
				if !requests.Memory().IsZero() && !limits.Memory().IsZero() && requests.Memory().Cmp(*limits.Memory()) > 0 {
					errs = append(errs, fmt.Errorf("deployment: %s, container: %s requests.memory must small than limits.memory", deployment.Name, c.Name))
				}

				if requests.Memory().IsZero() {
					errs = append(errs, fmt.Errorf("deployment: %s, container: %s must set memory request", deployment.Name, c.Name))
				} else {
					requiredMemory += requests.Memory().AsApproximateFloat64()
				}
				if requests.Cpu().IsZero() {
					errs = append(errs, fmt.Errorf("deployment: %s, container: %s must set cpu request", deployment.Name, c.Name))
				} else {
					requiredCPU += requests.Cpu().AsApproximateFloat64()
				}
				if limits.Memory().IsZero() {
					errs = append(errs, fmt.Errorf("deployment: %s, container: %s must set memory limit", deployment.Name, c.Name))
				} else {
					limitMemory += limits.Memory().AsApproximateFloat64()
				}
				if limits.Cpu().IsZero() {
					errs = append(errs, fmt.Errorf("deployment: %s, container: %s must set cpu limit", deployment.Name, c.Name))
				} else {
					limitCPU += limits.Cpu().AsApproximateFloat64()
				}
			}
		}
		if kind == StatefulSet {
			var sts v1.StatefulSet
			err := scheme.Scheme.Convert(r.Object, &sts, nil)
			if err != nil {
				return err
			}
			for _, c := range sts.Spec.Template.Spec.Containers {
				requests := c.Resources.Requests
				limits := c.Resources.Limits
				if !requests.Cpu().IsZero() && !limits.Cpu().IsZero() && requests.Cpu().Cmp(*limits.Cpu()) > 0 {
					errs = append(errs, fmt.Errorf("deployment: %s, container: %s requests.cpu must small than limits.cpu", sts.Name, c.Name))
				}
				if !requests.Memory().IsZero() && !limits.Memory().IsZero() && requests.Memory().Cmp(*limits.Memory()) > 0 {
					errs = append(errs, fmt.Errorf("deployment: %s, container: %s requests.memory must small than limits.memory", sts.Name, c.Name))
				}
				if requests.Memory().IsZero() {
					errs = append(errs, fmt.Errorf("statefulset: %s, container: %s must set memory request", sts.Name, c.Name))
				} else {
					requiredMemory += requests.Memory().AsApproximateFloat64()
				}
				if requests.Cpu().IsZero() {
					errs = append(errs, fmt.Errorf("statefulset: %s, container: %s must set cpu request", sts.Name, c.Name))
				} else {
					requiredCPU += requests.Cpu().AsApproximateFloat64()
				}
				if limits.Memory().IsZero() {
					errs = append(errs, fmt.Errorf("statefulset: %s, container: %s must set memory limit", sts.Name, c.Name))
				} else {
					limitMemory += limits.Memory().AsApproximateFloat64()
				}
				if limits.Cpu().IsZero() {
					errs = append(errs, fmt.Errorf("statefulset: %s, container: %s must set cpu limit", sts.Name, c.Name))
				} else {
					limitCPU += limits.Cpu().AsApproximateFloat64()
				}
			}
		}
	}
	if limitCPU > appLimitedCPU {
		errs = append(errs, fmt.Errorf("sum of all containers resources limits cpu should less than OlaresManifest.yaml spec.limitedCpu"))
	}
	if limitMemory > appLimitedMemory {
		errs = append(errs, fmt.Errorf("sum of all containers resources limits memory should less than OlaresManifest.yaml spec.limitedMemory"))
	}
	if requiredCPU > appRequiredCPU {
		errs = append(errs, fmt.Errorf("sum of all containers resources requests cpu should less than OlaresManifest.yaml spec.requiredCpu"))
	}
	if requiredMemory > appRequiredMemory {
		errs = append(errs, fmt.Errorf("sum of all containers resources requests memory should less than OlaresManifest.yaml spec.requiredMemory"))
	}
	return AggregateErr(errs)
}

func CheckResource(oacPath string, cfg *AppConfiguration, options *LintOptions) error {
	resources, err := getResourceListFromChart(oacPath, cfg, options)
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}
	err = checkResourceLimit(resources, cfg)
	if err != nil {
		return err
	}
	err = checkUploadConfig(resources, cfg)
	return err
}

func checkResourceNamespace(resources kube.ResourceList) error {
	errs := make([]error, 0)
	for _, r := range resources {
		kind := r.Object.GetObjectKind().GroupVersionKind().Kind
		if kind == Deployment || kind == StatefulSet || kind == DaemonSet {
			if r.Namespace != "app-namespace" {
				err := fmt.Errorf("illegal namespace: %s for %s, name %s", r.Namespace, kind, r.Name)
				errs = append(errs, err)
			}
		} else {
			if r.Namespace != "app-namespace" && !strings.HasPrefix(r.Namespace, "user-system-") {
				err := fmt.Errorf("illegal namespace: %s for %s, name %s", r.Namespace, kind, r.Name)
				errs = append(errs, err)
			}
		}
	}
	return AggregateErr(errs)
}

func checkUploadConfig(resources kube.ResourceList, cfg *AppConfiguration) error {
	if cfg.Options.Upload == nil {
		return nil
	}
	var err error
	for _, r := range resources {
		kind := r.Object.GetObjectKind().GroupVersionKind().Kind
		if kind == Deployment {
			var deployment v1.Deployment
			err = scheme.Scheme.Convert(r.Object, &deployment, nil)
			if err != nil {
				return err
			}
			for _, c := range deployment.Spec.Template.Spec.Containers {
				for _, v := range c.VolumeMounts {
					if filepath.Clean(v.MountPath) == filepath.Clean(cfg.Options.Upload.Dest) {
						return nil
					}
				}
			}
		}
		if kind == StatefulSet {
			var sts v1.StatefulSet
			err = scheme.Scheme.Convert(r.Object, &sts, nil)
			if err != nil {
				return err
			}
			err = scheme.Scheme.Convert(r.Object, &sts, nil)
			if err != nil {
				return err
			}
			for _, c := range sts.Spec.Template.Spec.Containers {
				for _, v := range c.VolumeMounts {
					if filepath.Clean(v.MountPath) == filepath.Clean(cfg.Options.Upload.Dest) {
						return nil
					}
				}
			}
		}
	}
	return fmt.Errorf("can not find volumemount path equal upload Dest: %s", cfg.Options.Upload.Dest)
}
