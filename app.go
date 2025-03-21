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
	ConfigVersion string      `yaml:"olaresManifest.version" json:"olaresManifest.version" vd:"len($)>0;msg:sprintf('invalid parameter: %v;olaresManifest.version must satisfy the expr: len($)>0',$)"`
	Metadata      AppMetaData `yaml:"metadata" json:"metadata"`
	Entrances     []Entrance  `yaml:"entrances" json:"entrances" vd:"len($)>0 && len($)<=10;msg:sprintf('invalid parameter: %v;entrances must satisfy the expr: len($)>0 && len($)<=10',$)"`
	Spec          AppSpec     `yaml:"spec,omitempty" json:"spec,omitempty"`
	Permission    *Permission `yaml:"permission" json:"permission" vd:"?"`
	Middleware    *Middleware `yaml:"middleware,omitempty" json:"middleware,omitempty" vd:"?"`
	Options       *Options    `yaml:"options" json:"options" vd:"?"`
}

type Middleware struct {
	Postgres   *PostgresConfig   `yaml:"postgres,omitempty" json:"postgres,omitempty" vd:"?"`
	Redis      *RedisConfig      `yaml:"redis,omitempty" json:"redis,omitempty" vd:"?"`
	MongoDB    *MongodbConfig    `yaml:"mongodb,omitempty" json:"mongodb,omitempty" vd:"?"`
	ZincSearch *ZincSearchConfig `yaml:"zincSearch,omitempty" json:"zincSearch,omitempty" vd:"?"`
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

type ZincSearchConfig struct {
	Username string  `yaml:"username" json:"username" vd:"len($)>0;msg:sprintf('invalid parameter: %v;username must satisfy the expr: len($)>0',$)"`
	Password string  `yaml:"password,omitempty" json:"password,omitempty" vd:"-"`
	Indexes  []Index `yaml:"indexes" json:"indexes" vd:"len($)>0;msg:sprintf('invalid parameter: %v;indexes must satisfy the expr: len($)>0',$)"`
}

type Index struct {
	Name string `yaml:"name" json:"name" vd:"len($)>0;msg:sprintf('invalid parameter: %v;name must satisfy the expr: len($)>0',$)"`
}
type AppMetaData struct {
	Name        string   `yaml:"name" json:"name"  vd:"len($)>0 && len($)<=30;msg:sprintf('invalid parameter: %v;name must satisfy the expr: len($)>0 && len($)<=30',$)"`
	Icon        string   `yaml:"icon" json:"icon" vd:"len($)>0;msg:sprintf('invalid parameter: %v;icon must satisfy the expr: len($)>0',$)"`
	Description string   `yaml:"description" json:"description" vd:"len($)>0;msg:sprintf('invalid parameter: %v;description must satisfy the expr: len($)>0',$)"`
	AppID       string   `yaml:"appid" json:"appid" vd:"-"`
	Title       string   `yaml:"title" json:"title" vd:"len($)>0 && len($)<=30;msg:sprintf('invalid parameter: %v;title must satisfy the expr: len($)>0 && len($)<=30',$)"`
	Version     string   `yaml:"version" json:"version" vd:"len($)>0;msg:sprintf('invalid parameter: %v;version must satisfy the expr: len($)>0',$)"`
	Categories  []string `yaml:"categories" json:"categories"`
	Rating      float32  `yaml:"rating" json:"rating" vd:"-"`
	Target      string   `yaml:"target" json:"target" vd:"-"`
}

type Entrance struct {
	Name      string `yaml:"name" json:"name" vd:"regexp('^([a-z0-9A-Z-]*)$') && len($)<=63;msg:sprintf('invalid parameter: %v;name must satisfy the expr: regexp(^([a-z0-9-]*)$)',$)"`
	Host      string `yaml:"host" json:"host" vd:"regexp('^([a-z]([-a-z0-9]*[a-z0-9]))$') && len($)<=63;msg:sprintf('invalid parameter: %v;host must satisfy the expr: regexp(^([a-z]([-a-z0-9]*[a-z0-9]))$)',$)"`
	Port      int32  `yaml:"port" json:"port" vd:"$>0;msg:sprintf('invalid parameter: %v;port must satisfy the expr: $>0',$)"`
	Icon      string `yaml:"icon" json:"icon"`
	Title     string `yaml:"title" json:"title" vd:"len($)>0 && len($)<=30 && regexp('^([a-z0-9A-Z-\\s]*)$');msg:sprintf('invalid parameter: %v;title must satisfy the expr: len($)>0 && len($)<=30 && regexp(^([a-z0-9A-Z-\\s]*)$)',$)"`
	AuthLevel string `yaml:"authLevel" json:"authLevel"`
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
	AppData bool         `yaml:"appData" json:"appData"`
	SysData []SysDataCfg `yaml:"sysData" json:"sysData"`
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
	Policies     *[]Policy     `yaml:"policies" json:"policies" vd:"?"`
	Analytics    *Analytics    `yaml:"analytics" json:"analytics" vd:"?"`
	Dependencies *[]Dependency `yaml:"dependencies" json:"dependencies" vd:"?"`
	Upload       *Upload       `yaml:"upload" json:"upload"`
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
