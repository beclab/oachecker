package oachecker

import (
	"fmt"
	vd "github.com/bytedance/go-tagexpr/v2/validator"
)

func init() {
	vd.SetErrorFactory(func(failPath, msg string) error {
		return fmt.Errorf(`"validation failed: %s","msg": "%s"`, failPath, msg)
	})
}

type AppConfiguration struct {
	ConfigVersion string        `yaml:"olaresManifest.version" json:"olaresManifest.version" vd:"len($)>0;msg:sprintf('invalid parameter: %v;olaresManifest.version must satisfy the expr: len($)>0',$)"`
	ConfigType    string        `yaml:"olaresManifest.type" json:"olaresManifest.type"`
	Metadata      AppMetaData   `yaml:"metadata" json:"metadata"`
	Entrances     []Entrance    `yaml:"entrances" json:"entrances" vd:"len($)>0 && len($)<=10;msg:sprintf('invalid parameter: %v;entrances must satisfy the expr: len($)>0 && len($)<=10',$)"`
	Ports         []ServicePort `yaml:"ports" json:"ports"`
	TailScale     TailScale     `yaml:"tailScale" json:"tailScale"`
	Spec          AppSpec       `yaml:"spec,omitempty" json:"spec,omitempty"`
	// TODO:hys add validate for permission field
	Permission Permission  `yaml:"permission" json:"permission" vd:"?"`
	Middleware *Middleware `yaml:"middleware,omitempty" json:"middleware,omitempty" vd:"?"`
	Options    Options     `yaml:"options" json:"options" vd:"?"`
	Provider   []Provider  `yaml:"provider,omitempty" json:"provider,omitempty" description:"app provider information"`
	Envs       []AppEnvVar `yaml:"envs,omitempty" json:"envs,omitempty"`
}

type Provider struct {
	Name     string   `yaml:"name" json:"name"`
	Entrance string   `yaml:"entrance" json:"entrance"`
	Paths    []string `yaml:"paths" json:"paths"`
	Verbs    []string `yaml:"verbs" json:"verbs"`
}

type AppEnvVar struct {
	EnvVarSpec    `json:",inline" yaml:",inline"`
	ApplyOnChange bool       `json:"applyOnChange,omitempty" yaml:"applyOnChange,omitempty"`
	ValueFrom     *ValueFrom `json:"valueFrom,omitempty" yaml:"valueFrom,omitempty"`
}

// ValueFrom defines a reference to an environment variable (UserEnv or SystemEnv)
type ValueFrom struct {
	EnvName string `json:"envName" validate:"required"`
	Status  string `json:"status,omitempty"`
}

type EnvVarSpec struct {
	EnvName     string `json:"envName" yaml:"envName" validate:"required"`
	Value       string `json:"value,omitempty" yaml:"value,omitempty"`
	Default     string `json:"default,omitempty" yaml:"default,omitempty"`
	Editable    bool   `json:"editable,omitempty" yaml:"editable,omitempty"`
	Type        string `json:"type,omitempty" yaml:"type,omitempty"`
	Required    bool   `json:"required,omitempty" yaml:"required,omitempty"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
}

type Middleware struct {
	Postgres      *PostgresConfig      `yaml:"postgres,omitempty" json:"postgres,omitempty" vd:"?"`
	Redis         *RedisConfig         `yaml:"redis,omitempty" json:"redis,omitempty" vd:"?"`
	MongoDB       *MongodbConfig       `yaml:"mongodb,omitempty" json:"mongodb,omitempty" vd:"?"`
	Nats          *NatsConfig          `yaml:"nats,omitempty"`
	Minio         *MinioConfig         `yaml:"minio,omitempty"`
	RabbitMQ      *RabbitMQConfig      `yaml:"rabbitmq,omitempty"`
	Elasticsearch *ElasticsearchConfig `yaml:"elasticsearch,omitempty"`
	MariaDB       *MariaDBConfig       `yaml:"mariadb,omitempty"`
	MySQL         *MySQLConfig         `yaml:"mysql,omitempty"`
	Argo          *ArgoConfig          `yaml:"argo,omitempty"`
}

type RabbitMQConfig struct {
	Username string  `yaml:"username" json:"username"`
	Password string  `yaml:"password" json:"password"`
	VHosts   []VHost `yaml:"vhosts" json:"vhosts"`
}

type VHost struct {
	Name string `json:"name"`
}

type ElasticsearchConfig struct {
	Username string  `yaml:"username" json:"username"`
	Password string  `yaml:"password" json:"password"`
	Indexes  []Index `yaml:"indexes" json:"indexes"`
}

type Index struct {
	Name string `json:"name"`
}
type ArgoConfig struct {
	Required bool `yaml:"required" json:"required"`
}

type MinioConfig struct {
	Username string   `yaml:"username" json:"username"`
	Password string   `yaml:"password" json:"password"`
	Buckets  []Bucket `yaml:"buckets" json:"buckets"`
}

type Bucket struct {
	Name string `json:"name"`
}
type NatsConfig struct {
	Username string    `yaml:"username" json:"username"`
	Password string    `yaml:"password,omitempty" json:"password,omitempty"`
	Subjects []Subject `yaml:"subjects" json:"subjects"`
	Refs     []Ref     `yaml:"refs" json:"refs"`
}
type Ref struct {
	AppName string `yaml:"appName" json:"appName"`
	// option for ref app in user-space-<>, user-system-<>, os-system
	AppNamespace string       `yaml:"appNamespace" json:"appNamespace"`
	Subjects     []RefSubject `yaml:"subjects" json:"subjects"`
}

type RefSubject struct {
	Name string   `yaml:"name" json:"name"`
	Perm []string `yaml:"perm" json:"perm"`
}

type Subject struct {
	Name string `yaml:"name" json:"name"`
	// Permissions indicates the permission that app can perform on this subject
	Permission Permission   `yaml:"permission" json:"permission"`
	Export     []Permission `yaml:"export" json:"export"`
}

// MariaDBConfig contains fields for mariadb config.
type MariaDBConfig struct {
	Username  string     `yaml:"username" json:"username"`
	Password  string     `yaml:"password,omitempty" json:"password"`
	Databases []Database `yaml:"databases" json:"databases"`
}

// MySQLConfig contains fields for mysql config.
type MySQLConfig struct {
	Username  string     `yaml:"username" json:"username"`
	Password  string     `yaml:"password,omitempty" json:"password"`
	Databases []Database `yaml:"databases" json:"databases"`
}

type Database struct {
	Name        string `yaml:"name" json:"name" vd:"len($)>0;msg:sprintf('invalid parameter: %v;name must satisfy the expr: len($)>0',$)"`
	Distributed bool   `yaml:"distributed,omitempty" json:"distributed,omitempty" vd:"-"`
}

type PostgresConfig struct {
	Username  string     `yaml:"username" json:"username" vd:"len($)>0;msg:sprintf('invalid parameter: %v;username must satisfy the expr: len($)>0',$)"`
	Password  string     `yaml:"password,omitempty" json:"password,omitempty" vd:"-"`
	Databases []Database `yaml:"databases" json:"databases" vd:"len($)>0;msg:sprintf('invalid parameter: %v;databases must satisfy the expr: len($)>0',$)"`
}

type RedisConfig struct {
	Password  string `yaml:"password,omitempty" json:"password,omitempty" vd:"-"`
	Namespace string `yaml:"namespace" json:"namespace" vd:"len($)>0;msg:sprintf('invalid parameter: %v;namespace must satisfy the expr: len($)>0',$)"`
}

type MongodbConfig struct {
	Username  string     `yaml:"username" json:"username" vd:"len($)>0;msg:sprintf('invalid parameter: %v;username must satisfy the expr: len($)>0',$)"`
	Password  string     `yaml:"password,omitempty" json:"password,omitempty" vd:"-"`
	Databases []Database `yaml:"databases" json:"databases" vd:"len($)>0;msg:sprintf('invalid parameter: %v;databases must satisfy the expr: len($)>0',$)"`
}

type AppMetaData struct {
	Name        string   `yaml:"name" json:"name"  vd:"len($)>0 && len($)<=30;msg:sprintf('invalid parameter: %v;name must satisfy the expr: len($)>0 && len($)<=30',$)"`
	Icon        string   `yaml:"icon" json:"icon" vd:"len($)>0;msg:sprintf('invalid parameter: %v;icon must satisfy the expr: len($)>0',$)"`
	Description string   `yaml:"description" json:"description" vd:"len($)>0;msg:sprintf('invalid parameter: %v;description must satisfy the expr: len($)>0',$)"`
	AppID       string   `yaml:"appid" json:"appid" vd:"-"`
	Title       string   `yaml:"title" json:"title" vd:"len($)>0 && len($)<=30;msg:sprintf('invalid parameter: %v;title must satisfy the expr: len($)>0 && len($)<=30',$)"`
	Version     string   `yaml:"version" json:"version" vd:"len($)>0;msg:sprintf('invalid parameter: %v;version must satisfy the expr: len($)>0',$)"`
	Categories  []string `yaml:"categories" json:"categories"`
	//Rating      float32  `yaml:"rating" json:"rating" vd:"-"`
	Target string `yaml:"target" json:"target" vd:"-"`
}

type Entrance struct {
	Name            string `yaml:"name" json:"name" vd:"regexp('^([a-z0-9A-Z-]*)$') && len($)<=63;msg:sprintf('invalid parameter: %v;name must satisfy the expr: regexp(^([a-z0-9-]*)$)',$)"`
	Host            string `yaml:"host" json:"host" vd:"regexp('^([a-z]([-a-z0-9]*[a-z0-9]))$') && len($)<=63;msg:sprintf('invalid parameter: %v;host must satisfy the expr: regexp(^([a-z]([-a-z0-9]*[a-z0-9]))$)',$)"`
	Port            int32  `yaml:"port" json:"port" vd:"$>0;msg:sprintf('invalid parameter: %v;port must satisfy the expr: $>0',$)"`
	Icon            string `yaml:"icon" json:"icon"`
	Title           string `yaml:"title" json:"title" vd:"len($)>0 && len($)<=30 && regexp('^([a-z0-9A-Z-\\s]*)$');msg:sprintf('invalid parameter: %v;title must satisfy the expr: len($)>0 && len($)<=30 && regexp(^([a-z0-9A-Z-\\s]*)$)',$)"`
	AuthLevel       string `yaml:"authLevel" json:"authLevel"`
	Invisible       bool   `yaml:"invisible,omitempty" json:"invisible,omitempty"`
	OpenMethod      string `yaml:"openMethod" json:"openMethod"`
	WindowPushState bool   `yaml:"windowPushState,omitempty" json:"windowPushState,omitempty"`
}

type AppSpec struct {
	VersionName        string         `yaml:"versionName,omitempty" json:"versionName"`
	FullDescription    string         `yaml:"fullDescription" json:"fullDescription"`
	UpgradeDescription string         `yaml:"upgradeDescription" json:"upgradeDescription"`
	PromoteImage       []string       `yaml:"promoteImage" json:"promoteImage"`
	PromoteVideo       string         `yaml:"promoteVideo" json:"promoteVideo"`
	SubCategory        string         `yaml:"subCategory" json:"subCategory"`
	Developer          string         `yaml:"developer" json:"developer"`
	RequiredMemory     string         `yaml:"requiredMemory" json:"requiredMemory" vd:"regexp('^(?:\\d+(?:\\.\\d+)?(?:[eE][-+]?(\\d+|i))?(?:[kKMGTP]?i?|[mMGTPE])?|[kKMGTP]i|[mMGTPE])$');msg:sprintf('invalid parameter: %v;requiredMemory must satisfy the expr: regexp(^(?:\\d+(?:\\.\\d+)?(?:[eE][-+]?(\\d+|i))?(?:[kKMGTP]?i?|[mMGTPE])?|[kKMGTP]i|[mMGTPE])$)',$)"`
	RequiredDisk       string         `yaml:"requiredDisk" json:"requiredDisk"     vd:"regexp('^(?:\\d+(?:\\.\\d+)?(?:[eE][-+]?(\\d+|i))?(?:[kKMGTP]?i?|[mMGTPE])?|[kKMGTP]i|[mMGTPE])$');msg:sprintf('invalid parameter: %v;requiredDisk must satisfy the expr: regexp(^(?:\\d+(?:\\.\\d+)?(?:[eE][-+]?(\\d+|i))?(?:[kKMGTP]?i?|[mMGTPE])?|[kKMGTP]i|[mMGTPE])$)',$)"`
	SupportClient      *SupportClient `yaml:"supportClient" json:"supportClient" vd:"?"`
	SupportArch        []string       `yaml:"supportArch" json:"supportArch"`
	RequiredGPU        string         `yaml:"requiredGpu" json:"requiredGpu" vd:"len($)==0 || regexp('^(?:\\d+(?:\\.\\d+)?(?:[eE][-+]?(\\d+|i))?(?:[kKMGTP]?i?|[mMGTPE])?|[kKMGTP]i|[mMGTPE])$');msg:sprintf('invalid parameter: %v;requiredGpu must satisfy the expr: len($) == 0 || regexp(^(?:\\d+(?:\\.\\d+)?(?:[eE][-+]?(\\d+|i))?(?:[kKMGTP]?i?|[mMGTPE])?|[kKMGTP]i|[mMGTPE])$)',$)"`
	RequiredCPU        string         `yaml:"requiredCpu" json:"requiredCpu" vd:"regexp('^(?:\\d+(?:\\.\\d+)?(?:[eE][-+]?(\\d+|i))?(?:[kKMGTP]?i?|[mMGTPE])?|[kKMGTP]i|[mMGTPE])$');msg:sprintf('invalid parameter: %v;requiredCpu must satisfy the expr: regexp(^(?:\\d+(?:\\.\\d+)?(?:[eE][-+]?(\\d+|i))?(?:[kKMGTP]?i?|[mMGTPE])?|[kKMGTP]i|[mMGTPE])$)',$)"`
	LimitedMemory      string         `yaml:"limitedMemory" json:"limitedMemory" vd:"regexp('^(?:\\d+(?:\\.\\d+)?(?:[eE][-+]?(\\d+|i))?(?:[kKMGTP]?i?|[mMGTPE])?|[kKMGTP]i|[mMGTPE])$');msg:sprintf('invalid parameter: %v;limitedMemory must satisfy the expr: regexp(^(?:\\d+(?:\\.\\d+)?(?:[eE][-+]?(\\d+|i))?(?:[kKMGTP]?i?|[mMGTPE])?|[kKMGTP]i|[mMGTPE])$)',$)"`
	LimitedCPU         string         `yaml:"limitedCpu" json:"limitedCpu" vd:"regexp('^(?:\\d+(?:\\.\\d+)?(?:[eE][-+]?(\\d+|i))?(?:[kKMGTP]?i?|[mMGTPE])?|[kKMGTP]i|[mMGTPE])$');msg:sprintf('invalid parameter: %v;limitedCpu must satisfy the expr: regexp(^(?:\\d+(?:\\.\\d+)?(?:[eE][-+]?(\\d+|i))?(?:[kKMGTP]?i?|[mMGTPE])?|[kKMGTP]i|[mMGTPE])$)',$)"`

	RunAsUser           bool      `yaml:"runAsUser" json:"runAsUser"`
	RunAsInternal       bool      `yaml:"runAsInternal" json:"runAsInternal"`
	PodGPUConsumePolicy string    `yaml:"podGpuConsumePolicy" json:"podGpuConsumePolicy"`
	SubCharts           []ChartV2 `yaml:"subCharts" json:"subCharts"`

	Language     []string     `yaml:"language,omitempty" json:"language,omitempty"`
	Submitter    string       `yaml:"submitter,omitempty" json:"submitter,omitempty"`
	Doc          string       `yaml:"doc,omitempty" json:"doc,omitempty"`
	Website      string       `yaml:"website,omitempty" json:"website,omitempty"`
	FeatureImage string       `yaml:"featuredImage,omitempty" json:"featuredImage,omitempty"`
	SourceCode   string       `yaml:"sourceCode,omitempty" json:"sourceCode,omitempty"`
	License      []TextAndURL `yaml:"license,omitempty" json:"license,omitempty"`
	Legal        []TextAndURL `yaml:"legal,omitempty" json:"legal,omitempty"`
}

type TextAndURL struct {
	Text string `yaml:"text" json:"text" bson:"text"`
	URL  string `yaml:"url" json:"url" bson:"url"`
}

type SupportClient struct {
	Edge    string `yaml:"edge" json:"edge"`
	Android string `yaml:"android" json:"android"`
	Ios     string `yaml:"ios" json:"ios"`
	Windows string `yaml:"windows" json:"windows"`
	Mac     string `yaml:"mac" json:"mac"`
	Linux   string `yaml:"linux" json:"linux"`
}

type Permission struct {
	AppData        bool                 `yaml:"appData,omitempty" json:"appData,omitempty"  description:"app data permission for writing"`
	AppCache       bool                 `yaml:"appCache" json:"appCache"`
	UserData       []string             `yaml:"userData" json:"userData"`
	Provider       []ProviderPermission `yaml:"provider" json:"provider"  description:"system shared data permission for accessing"`
	ServiceAccount *string              `yaml:"serviceAccount,omitempty" json:"serviceAccount,omitempty" description:"service account for app permission"`
}

type ProviderPermission struct {
	AppName      string `yaml:"appName" json:"appName"`
	Namespace    string `yaml:"namespace,omitempty" json:"namespace,omitempty"`
	ProviderName string `yaml:"providerName" json:"providerName"`
}

type SysDataCfg struct {
	Group    string   `yaml:"group" json:"group" vd:"len($)>0;msg:sprintf('invalid parameter: %v;group must satisfy the expr: len($)>0',$)"`
	DataType string   `yaml:"dataType" json:"dataType" vd:"len($)>0;msg:sprintf('invalid parameter: %v;dataType must satisfy the expr: len($)>0',$)"`
	Version  string   `yaml:"version" json:"version" vd:"len($)>0;msg:sprintf('invalid parameter: %v;version must satisfy the expr: len($)>0',$)"`
	Ops      []string `yaml:"ops" json:"ops" vd:"len($)>0;msg:sprintf('invalid parameter: %v;ops must satisfy the expr: len($)>0',$)"`
}

type Policy struct {
	Description string `yaml:"description" json:"description" vd:"-"`
	URIRegex    string `yaml:"uriRegex" json:"uriRegex" vd:"len($)>0;msg:sprintf('invalid parameter: %v;uriRegex must satisfy the expr: len($)>0',$)"`
	Level       string `yaml:"level" json:"level" vd:"len($)>0;msg:sprintf('invalid parameter: %v;level must satisfy the expr: len($)>0',$)"`
	OneTime     bool   `yaml:"oneTime" json:"oneTime"`
	Duration    string `yaml:"validDuration" json:"validDuration" vd:"len($)==0 ||regexp('^((?:[-+]?\\d+(?:\\.\\d+)?([smhdwy]|us|ns|ms))+)$');msg:sprintf('invalid parameter: %v;validDuration must satisfy the expr: regexp(^((?:[-+]?\\d+(?:\\.\\d+)?([smhdwy]|us|ns|ms))+)$)',$)"`
}

type Options struct {
	MobileSupported      bool                     `yaml:"mobileSupported" json:"mobileSupported"`
	Policies             *[]Policy                `yaml:"policies" json:"policies" vd:"?"`
	Analytics            *Analytics               `yaml:"analytics" json:"analytics" vd:"?"`
	ResetCookie          *ResetCookie             `yaml:"resetCookie" json:"resetCookie" vd:"?"`
	Dependencies         *[]Dependency            `yaml:"dependencies" json:"dependencies" vd:"?"`
	AppScope             *AppScope                `yaml:"appScope" json:"appScope" vd:"?"`
	WsConfig             *WsConfig                `yaml:"wsConfig" json:"wsConfig" vd:"?"`
	Upload               *Upload                  `yaml:"upload" json:"upload"`
	SyncProvider         []map[string]interface{} `yaml:"syncProvider" json:"syncProvider"`
	OIDC                 OIDC                     `yaml:"oidc" json:"oidc"`
	ApiTimeout           *int64                   `yaml:"apiTimeout" json:"apiTimeout"`
	AllowedOutboundPorts []int                    `yaml:"allowedOutboundPorts" json:"AllowedOutboundPorts"`
	Images               []string                 `yaml:"images" json:"images"`
}

type OIDC struct {
	Enabled      bool   `yaml:"enabled" json:"enabled"`
	RedirectUri  string `yaml:"redirectUri" json:"redirectUri"`
	EntranceName string `yaml:"entranceName" json:"entranceName"`
}

type WsConfig struct {
	Port int    `yaml:"port" json:"port"`
	URL  string `yaml:"url" json:"url"`
}
type AppScope struct {
	ClusterScoped bool     `yaml:"clusterScoped" json:"clusterScoped"`
	AppRef        []string `yaml:"appRef" json:"appRef"`
}

type ResetCookie struct {
	Enabled bool `yaml:"enabled" json:"enabled"`
}

type Analytics struct {
	Enabled bool `yaml:"enabled" json:"enabled"`
}

type Dependency struct {
	Name    string `yaml:"name" json:"name" vd:"len($)>0;msg:sprintf('invalid parameter: %v;name must satisfy the expr: len($)>0',$)"`
	Version string `yaml:"version" json:"version" vd:"len($)>0;msg:sprintf('invalid parameter: %v;version must satisfy the expr: len($)>0',$)"`
	// dependency type: system, application.
	Type string `yaml:"type" json:"type" vd:"$=='system' || $=='application';msg:sprintf('invalid parameter: %v;type must satisfy the expr: $==system || $==application',$)"`
}

type Upload struct {
	FileType    []string `yaml:"fileType" json:"fileType"`
	Dest        string   `yaml:"dest" json:"dest"`
	LimitedSize int      `yaml:"limitedSize" json:"limitedSize"`
}

type Mappings struct {
	Properties map[string]Property `json:"properties,omitempty"`
}

type Property struct {
	Type           string `json:"type"` // text, keyword, date, numeric, boolean, geo_point
	Analyzer       string `json:"analyzer,omitempty"`
	SearchAnalyzer string `json:"search_analyzer,omitempty"`
	Format         string `json:"format,omitempty"`    // date format yyyy-MM-dd HH:mm:ss || yyyy-MM-dd || epoch_millis
	TimeZone       string `json:"time_zone,omitempty"` // date format time_zone
	Index          bool   `json:"index"`
	Store          bool   `json:"store"`
	Sortable       bool   `json:"sortable"`
	Aggregatable   bool   `json:"aggregatable"`
	Highlightable  bool   `json:"highlightable"`

	Fields map[string]Property `json:"fields,omitempty"`
}

// Chart represents the structure of the Chart.yaml file
type Chart struct {
	APIVersion string `yaml:"apiVersion"`
	Name       string `yaml:"name"`
	Version    string `yaml:"version"`
}

type TitleInfo struct {
	PrType  string
	Folder  string
	Version string
}

type ServicePort struct {
	Name string `json:"name" yaml:"name"`
	Host string `yaml:"host" json:"host"`
	Port int32  `yaml:"port" json:"port"`

	ExposePort int32 `yaml:"exposePort" json:"exposePort,omitempty"`

	// The protocol for this entrance. Supports "tcp" and "udp","".
	// Default is tcp/udp, "" mean tcp and udp.
	// +default="tcp/udp"
	// +optional
	Protocol          string `yaml:"protocol" json:"protocol,omitempty"`
	AddToTailscaleAcl bool   `yaml:"addToTailscaleAcl" json:"addToTailscaleAcl"`
}

type ACL struct {
	Action string   `json:"action,omitempty"`
	Src    []string `json:"src,omitempty"`
	Proto  string   `json:"proto"`
	Dst    []string `json:"dst"`
}

type TailScale struct {
	ACLs      []ACL    `json:"acls,omitempty"`
	SubRoutes []string `json:"subRoutes,omitempty"`
}

type ChartV2 struct {
	Name   string `yaml:"name" json:"name"`
	Shared bool   `yaml:"shared" json:"shared"`
}
