package oachecker

const (
	RoleBinding        = "RoleBinding"
	Role               = "Role"
	ClusterRole        = "ClusterRole"
	ServiceAccount     = "ServiceAccount"
	ClusterRoleBinding = "ClusterRoleBinding"
	Deployment         = "Deployment"
	StatefulSet        = "StatefulSet"
	DaemonSet          = "DaemonSet"
	ManifestName       = "OlaresManifest.yaml"
	ManifestRenderKey  = "chart/OlaresManifest.yaml"
)

const RULES = `rules:
- apiGroups:
  - '*'
  resources:
  - nodes
  - networkpolicies
  verbs:
  - create
  - update
  - patch
  - delete
  - deletecollection
`

var (
	ReOpenMsg                 = "Please create a new PR instead of reopening this PR."
	TitleInvalid              = "Invalid PR format. PR title must conform to the following format: [pr type][foldername][version]title"
	PrOwnerInvalid            = "Authorization exceptions. [%s] owner invalid"
	PrAlreadyExist            = "There is already an open PR for this folder %d"
	PrFolderExist             = "[%s] already exists. Please check your folder name or use update PR to modify."
	PrSpecialFiles            = "Invalid change. There should be no special control files in your submission."
	PrTypeChange              = "Please recreate a PR, if you want to modify the PR type."
	PrPermission              = "Authorization exceptions. You do not have permission to modify [%s]."
	PrVersionMustIncrease     = "Invalid change. Update version %s need to be greater than %s from main branch."
	PrDraftPass               = "After the modification is completed, please click ready for review to submit the PR"
	PrConflict                = "Conflict needs to be resolved"
	PrNotExist                = "[%s] does not exist. Please check your folder name or use new PR to submit a new one."
	PrShouldOnlyIncludeRemove = `Invalid change. Submissions should only include a ".remove" file`
	PrSuspendShould           = `Invalid change. Submissions should contain a ".suspend" file and no other special control files`
	PrCheckPass               = "Check passed, please wait for auto-merge."
	PrCheckPassLabel          = "Check skipped. PR has label [%s]. Please wait for auto-merge."
	PrFolderDif               = "Inconsistent info. Changed folder:%s is different from the foldername in title:%s"
	PrMultiDir                = "Invalid change. Change in multiple directory detected: %v. You should only modify one directory at a time."

	PrReviewerChangesRequested = "1, reviewer changes requested"
	PrReviewerApproved         = "4, reviewer approved"

	FolderNoOwners     = "Authorization exceptions. [%s] has no owners"
	CloneCodeFailed    = "clone code failed err [%s]"
	CheckoutCodeFailed = "checkout code failed err [%s]"
	ChartInvalid       = "Chart invalid. Error message: [%s]"

	AppCfgInfoEmtpyFromMain = "The info of OlaresManifest.yaml from the main branch is empty."
	AppCfgReadFailed        = "Failed to read OlaresManifest.yaml in folder [%s]: [%v]"
	AppCfgParseFailed       = "Failed to parse OlaresManifest.yaml in folder [%s]: [%v]"

	DbErr = "db err[%s]"

	//chart invalid
	InvalidAppCfgType            = "olaresManifest.type %s invalid, must in %v"
	InvalidFolderName            = "invalid folder name: '%s' must '^[a-z0-9]{1,30}$'"
	FolderNotExist               = "folder does not exist: '%s'"
	MissingChartYaml             = "missing Chart.yaml in folder: '%s'"
	ReadChartYamlFailed          = "failed to read Chart.yaml in folder '%s': %v"
	ParseChartYamlFailed         = "failed to parse Chart.yaml in folder '%s': %v"
	ApiVersionFieldEmptyInAppCfg = "apiVersion field empty in OlaresManifest.yaml in chart '%s'"
	NameFieldEmptyInAppCfg       = "name field empty in OlaresManifest.yaml in chart '%s'"
	VersionFieldEmptyInAppCfg    = "version field empty in OlaresManifest.yaml in chart '%s'"
	MissingValuesYaml            = "missing values.yaml in folder: '%s'"
	MissingTemplatesFolder       = "missing templates folder in folder: '%s'"
	MissingAppCfg                = "missing OlaresManifest.yaml in folder: '%s'"
	ReadAppCfgFailed             = "failed to read OlaresManifest.yaml in folder '%s': %v"
	ParseAppCfgFailed            = "failed to parse OlaresManifest.yaml in folder '%s': %v"
	NameMustSame1                = "inconsistent info. name must be the same in chart. name in Chart.yaml:%s, chartFolder:%s, OlaresManifest.yaml:%s"
	NameMustSame2                = "inconsistent info. name must be the same in chart. name in Chart.yaml:%s, chartFolder:%s, folder in title:%s, OlaresManifest.yaml:%s"
	VersionMustSame1             = "inconsistent info. Version must be the same in chart. version in OlaresManifest.yaml:%s, Chart.yaml:%s"
	VersionMustSame2             = "inconsistent info. Version must be the same in chart. version in OlaresManifest.yaml:%s, Chart.yaml:%s, title:%s"
	InvalidCategories            = "categories %v invalid, must in %v"
	FolderNameInvalid            = "foldername %s in reserved foldername list, invalid"
	//chart image invalid
	InvalidPromoteImageFormat = "invalid promote image format: %s, only support png, jpeg, webp"
	InvalidIconImageFormat    = "invalid icon image format: %s, only support png, webp"
	ImageNotFound             = "info not found in image"
	ImageSizeExceed           = "image size %d exceeds the limit: %d"
	InvalidImageDimensions    = "invalid image dimensions: %dx%d, should be: %dx%d"
)
