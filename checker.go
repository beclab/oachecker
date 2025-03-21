package oachecker

import vd "github.com/bytedance/go-tagexpr/v2/validator"

type LintOptions struct {
	Owner             string
	Admin             string
	SkipManifestCheck bool
	SkipResourceCheck bool

	SkipFolderCheck  bool
	CustomValidators []func(string, *AppConfiguration) error
}

func DefaultLintOptions() *LintOptions {
	return &LintOptions{
		Owner:             "",
		Admin:             "",
		SkipManifestCheck: false,
		SkipResourceCheck: false,
		SkipFolderCheck:   false,
		CustomValidators:  []func(string, *AppConfiguration) error{},
	}
}

func (o *LintOptions) WithOwner(owner string) *LintOptions {
	o.Owner = owner
	return o
}

func (o *LintOptions) WithAdmin(admin string) *LintOptions {
	o.Admin = admin
	return o
}

func (o *LintOptions) WithSameOwnerAndAdmin(value string) *LintOptions {
	o.Owner = value
	o.Admin = value
	return o
}

func (o *LintOptions) WithCustomValidator(validator func(string, *AppConfiguration) error) *LintOptions {
	o.CustomValidators = append(o.CustomValidators, validator)
	return o
}

func (o *LintOptions) WithAppDataValidator() {
	o.CustomValidators = append(o.CustomValidators, CheckAppData)
}

func (o *LintOptions) SkipManifest() *LintOptions {
	o.SkipManifestCheck = true
	return o
}

func (o *LintOptions) SkipResources() *LintOptions {
	o.SkipResourceCheck = true
	return o
}

func CheckChart(oacPath string) (err error) {
	err = CheckChartFolder(oacPath)
	if err != nil {
		return err
	}
	err = CheckAppCfg(oacPath)
	if err != nil {
		return err
	}
	err = CheckServiceAccountRole(oacPath)
	if err != nil {
		return err
	}
	return nil
}

func CheckManifest(oacPath string, cfg *AppConfiguration) error {
	err := vd.Validate(cfg, true)
	if err != nil {
		return err
	}
	err = CheckSupportedArch(cfg)
	if err != nil {
		return err
	}

	err = CheckAppEntrances(cfg)
	if err != nil {
		return err
	}
	return nil
}

func CheckManifestFromFile(oacPath string, opts ...func(map[string]interface{})) error {
	cfg, err := GetAppConfiguration(oacPath, opts...)
	if err != nil {
		return err
	}
	err = vd.Validate(cfg, true)
	if err != nil {
		return err
	}
	err = CheckSupportedArch(cfg)
	if err != nil {
		return err
	}

	err = CheckAppEntrances(cfg)
	if err != nil {
		return err
	}
	return nil
}

func CheckManifestFromContent(content []byte, opts ...func(map[string]interface{})) error {
	cfg, err := GetAppConfigurationFromContent(content, opts...)
	if err != nil {
		return err
	}
	err = vd.Validate(cfg, true)
	if err != nil {
		return err
	}
	err = CheckSupportedArch(cfg)
	if err != nil {
		return err
	}

	err = CheckAppEntrances(cfg)
	if err != nil {
		return err
	}
	return nil
}

func Lint(oacPath string, options *LintOptions) error {
	if options == nil {
		options = DefaultLintOptions()
	}

	var opts []func(map[string]interface{})
	if options.Owner != "" {
		opts = append(opts, WithOwner(options.Owner))
	}
	if options.Admin != "" {
		opts = append(opts, WithAdmin(options.Admin))
	}

	cfg, err := GetAppConfiguration(oacPath, opts...)
	if err != nil {
		return err
	}

	if !options.SkipManifestCheck {
		err = CheckManifest(oacPath, cfg)
		if err != nil {
			return err
		}
	}

	for _, validator := range options.CustomValidators {
		if err := validator(oacPath, cfg); err != nil {
			return err
		}
	}

	if !options.SkipResourceCheck {
		err = CheckResource(oacPath, cfg, options)
		if err != nil {
			return err
		}
	}

	if !options.SkipFolderCheck {
		err = CheckChartFolder(oacPath)
		if err != nil {
			return err
		}
	}

	return nil
}

func LintWithDefaultOptions(oacPath string) error {
	return Lint(oacPath, DefaultLintOptions())
}

func LintWithSameOwnerAdmin(oacPath string, ownerAdmin string) error {
	options := DefaultLintOptions().WithSameOwnerAndAdmin(ownerAdmin)
	return Lint(oacPath, options)
}

func LintWithDifferentOwnerAdmin(oacPath string, owner string, admin string) error {
	options := DefaultLintOptions().WithOwner(owner).WithAdmin(admin)
	return Lint(oacPath, options)
}
