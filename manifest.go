package oachecker

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	vd "github.com/bytedance/go-tagexpr/v2/validator"
	"gopkg.in/yaml.v3"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/engine"
	"io"
	"k8s.io/apimachinery/pkg/util/sets"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	defaultOwner = "default"
	defaultAdmin = "default"
)

func getAppConfigFromCfg(f io.ReadCloser, opts ...func(map[string]interface{})) (*AppConfiguration, error) {
	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	renderedData, err := RenderManifestFromContent(data, opts...)
	if err != nil {
		return nil, err
	}
	var cfg AppConfiguration
	if err := yaml.Unmarshal([]byte(renderedData), &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

//func getAppConfigFromCfgFile(oacPath string, owner, admin string) (*AppConfiguration, error) {
//	if !strings.HasSuffix(oacPath, "/") {
//		oacPath += "/"
//	}
//	f, err := os.Open(oacPath + "OlaresManifest.yaml")
//	if err != nil {
//		return nil, err
//	}
//	return getAppConfigFromCfg(f, owner, admin)
//}

func CheckAppCfg(oacPath string, opts ...func(map[string]interface{})) error {
	cfg, err := GetAppConfiguration(oacPath, opts...)
	if err != nil {
		return err
	}
	return checkAppCfg(cfg, oacPath, true)
}

func checkAppCfg(cfg *AppConfiguration, oacPath string, checkAll ...bool) error {
	err := vd.Validate(cfg, checkAll...)
	if err != nil {
		return err
	}
	//err = CheckAppEntrances(cfg)
	//if err != nil {
	//	return err
	//}
	err = CheckSupportedArch(cfg)
	if err != nil {
		return err
	}

	err = CheckAppData(oacPath, cfg)
	if err != nil {
		return err
	}
	return CheckResource(oacPath, cfg, nil)
}

func CheckSupportedArch(cfg *AppConfiguration) error {
	if len(cfg.Spec.SupportArch) == 0 {
		return errors.New("spec.SupportArch can not be empty")
	}
	allSupportedArch := sets.String{"amd64": sets.Empty{}, "arm32v5": sets.Empty{}, "arm32v6": sets.Empty{},
		"arm32v7": sets.Empty{}, "arm64v8": sets.Empty{}, "i386": sets.Empty{}, "ppc64le": sets.Empty{},
		"s390x": sets.Empty{}, "mips64le": sets.Empty{}, "riscv64": sets.Empty{}, "windows-amd64": sets.Empty{}, "arm64": sets.Empty{}}
	for _, arch := range cfg.Spec.SupportArch {
		if !allSupportedArch.Has(arch) {
			return fmt.Errorf("unsupport arch: %s", arch)
		}
	}
	return nil
}

func CheckAppEntrances(cfg *AppConfiguration) error {
	//setsEntrance := sets.String{}
	setsName := sets.String{}
	for i, e := range cfg.Entrances {
		//entrance := fmt.Sprintf("%s:%d", e.Host, e.Port)
		//if setsEntrance.Has(entrance) {
		//	return fmt.Errorf("entrances:[%d] has replicated(entrance with same name and port were treat as same)", i)
		//}
		//setsEntrance.Insert(entrance)

		if setsName.Has(e.Name) {
			return fmt.Errorf("entrances:[%d] name has replicated", i)
		}
		setsName.Insert(e.Name)
	}
	return nil
}

func CheckAppData(oacPath string, cfg *AppConfiguration) error {
	if cfg.Permission.AppData {
		return nil
	}
	if !strings.HasSuffix(oacPath, "/") {
		oacPath += "/"
	}
	oacPath += "templates"
	p, err := regexp.Compile(`\.Values\.userspace\.appdata`)
	if err != nil {
		return err
	}
	var rerr error
	err = filepath.Walk(oacPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(path, ".yaml") {
			f, e := os.Open(path)
			if e != nil {
				return e
			}
			defer f.Close()
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				if p.MatchString(scanner.Text()) {
					rerr = fmt.Errorf("found .Values.userspace.appdata in %s, but not set permission.appData in OlaresManifest.yaml", filepath.Base(path))
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return rerr
}

func RenderManifestFromContent(content []byte, opts ...func(map[string]interface{})) (string, error) {
	values := map[string]interface{}{
		"admin": defaultAdmin,
		"bfl": map[string]string{
			"username": defaultOwner,
		},
	}

	for _, opt := range opts {
		opt(values)
	}
	c := &chart.Chart{
		Metadata: &chart.Metadata{
			Name:    "chart",
			Version: "0.0.1",
		},
		Templates: []*chart.File{
			{
				Name: ManifestName,
				Data: content,
			},
		},
	}

	valuesToRender, err := chartutil.ToRenderValues(c, values, chartutil.ReleaseOptions{}, nil)
	if err != nil {
		return "", err
	}

	e := engine.Engine{}
	renderedTemplates, err := e.Render(c, valuesToRender)
	if err != nil {
		return "", err
	}

	renderedYAML := renderedTemplates[ManifestRenderKey]

	return renderedYAML, nil
}

func WithOwner(owner string) func(map[string]interface{}) {
	return func(values map[string]interface{}) {
		if owner != "" {
			bfl := map[string]string{
				"username": owner,
			}
			values["bfl"] = bfl
		}
	}
}

func WithAdmin(admin string) func(map[string]interface{}) {
	return func(values map[string]interface{}) {
		if admin != "" {
			values["admin"] = admin
		}
	}
}

func GetAppConfiguration(oacPath string, opts ...func(map[string]interface{})) (*AppConfiguration, error) {
	if !strings.HasSuffix(oacPath, "/") {
		oacPath += "/"
	}
	f, err := os.Open(oacPath + ManifestName)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return getAppConfigFromCfg(f, opts...)
}

func GetAppConfigurationFromContent(content []byte, opts ...func(map[string]interface{})) (*AppConfiguration, error) {
	f := io.NopCloser(bytes.NewReader(content))
	return getAppConfigFromCfg(f, opts...)
}
