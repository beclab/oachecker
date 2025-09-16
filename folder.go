package oachecker

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

func baseChartFolderCheck(folder string) (*Chart, *AppConfiguration, string, error) {
	folderName := path.Base(folder)
	if !isValidFolderName(folderName) {
		return nil, nil, "", fmt.Errorf(InvalidFolderName, folder)
	}

	if !dirExists(folder) {
		return nil, nil, "", fmt.Errorf(FolderNotExist, folder)
	}

	chartFile := filepath.Join(folder, "Chart.yaml")
	if !fileExists(chartFile) {
		return nil, nil, "", fmt.Errorf(MissingChartYaml, folder)
	}

	chartContent, err := os.ReadFile(chartFile)
	if err != nil {
		return nil, nil, "", fmt.Errorf(ReadChartYamlFailed, folder, err)
	}
	var chart Chart
	if err := yaml.Unmarshal(chartContent, &chart); err != nil {
		return nil, nil, "", fmt.Errorf(ParseChartYamlFailed, folder, err)
	}

	if err := isValidChartFields(chart); err != nil {
		return nil, nil, "", err
	}

	valuesFile := filepath.Join(folder, "values.yaml")
	if !fileExists(valuesFile) {
		return nil, nil, "", fmt.Errorf(MissingValuesYaml, folder)
	}

	templatesDir := filepath.Join(folder, "templates")
	if !dirExists(templatesDir) {
		return nil, nil, "", fmt.Errorf(MissingTemplatesFolder, folder)
	}

	appCfgFile := filepath.Join(folder, "OlaresManifest.yaml")
	if !fileExists(appCfgFile) {
		return nil, nil, "", fmt.Errorf(MissingAppCfg, folder)
	}

	//appCfgContent, err := os.ReadFile(appCfgFile)
	//if err != nil {
	//	return nil, nil, "", fmt.Errorf(ReadAppCfgFailed, folder, err)
	//}
	//
	//var appConf AppConfiguration
	//if err := yaml.Unmarshal(appCfgContent, &appConf); err != nil {
	//	return nil, nil, "", fmt.Errorf(ParseAppCfgFailed, folder, err)
	//}
	appConf, err := GetAppConfiguration(folder)
	if err != nil {
		return nil, nil, "", fmt.Errorf(ReadAppCfgFailed, folder, err)
	}

	return &chart, appConf, folderName, nil
}

func CheckChartFolder(folder string) error { // todo extract func
	_, _, _, err := baseChartFolderCheck(folder)
	if err != nil {
		return err
	}

	return nil
}

func CheckSameVersion(folder string) error {
	chart, appConf, folderName, err := baseChartFolderCheck(folder)
	if err != nil {
		return err
	}

	err = isValidMetadataFields(appConf.Metadata, chart, folderName)
	if err != nil {
		return err
	}
	return nil
}

func CheckChartFolderWithTitle(folder string, titleInfo TitleInfo) error {
	chart, appConf, folderName, err := baseChartFolderCheck(folder)
	if err != nil {
		return err
	}

	if err = isValidMetadataFieldsWithTitle(appConf.Metadata, chart, folderName, titleInfo); err != nil {
		return err
	}

	if !checkCategories(appConf.Metadata.Categories) {
		return fmt.Errorf(InvalidCategories, appConf.Metadata.Categories, validCategoriesSlice)
	}

	if checkReservedWord(folderName) {
		return fmt.Errorf(FolderNameInvalid, folderName)
	}

	if err = CheckAppConfigImages(appConf); err != nil {
		return err
	}

	return nil
}

func checkReservedWord(str string) bool {
	reservedWords := []string{
		"user", "system", "space", "default", "os", "kubesphere", "kube",
		"kubekey", "kubernetes", "gpu", "tapr", "bfl", "bytetrade",
		"project", "pod",
	}

	for _, word := range reservedWords {
		if strings.EqualFold(str, word) {
			return true
		}
	}

	return false
}

func isValidFolderName(name string) bool {
	match, _ := regexp.MatchString("^[a-z0-9]{1,30}$", name)
	return match
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return (err == nil || os.IsExist(err)) && !info.IsDir()
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return (err == nil || os.IsExist(err)) && info.IsDir()
}

func isValidChartFields(chart Chart) error {
	if chart.APIVersion == "" {
		return fmt.Errorf(ApiVersionFieldEmptyInAppCfg, chart)
	}

	if chart.Name == "" {
		return fmt.Errorf(NameFieldEmptyInAppCfg, chart)
	}

	if chart.Version == "" {
		return fmt.Errorf(VersionFieldEmptyInAppCfg, chart)
	}

	return nil
}

func isValidMetadataFieldsWithTitle(metadata AppMetaData, chart *Chart, folder string, titleInfo TitleInfo) error {
	if chart.Name != folder || titleInfo.Folder != folder || metadata.Name != folder {
		return fmt.Errorf(NameMustSame2,
			chart.Name, folder, titleInfo.Folder, metadata.Name)
	}

	if metadata.Version != chart.Version || titleInfo.Version != chart.Version {
		return fmt.Errorf(VersionMustSame2, metadata.Version, chart.Version, titleInfo.Version)
	}

	return nil
}

func isValidMetadataFields(metadata AppMetaData, chart *Chart, folder string) error {
	if chart.Name != folder || metadata.Name != folder {
		return fmt.Errorf(NameMustSame1,
			chart.Name, folder, metadata.Name)
	}

	if metadata.Version != chart.Version {
		return fmt.Errorf(VersionMustSame1, metadata.Version, chart.Version)
	}

	return nil
}
