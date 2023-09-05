package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	changedCharts string
)

// Chart represents the structure of the Chart.yaml file
type Chart struct {
	APIVersion string `yaml:"apiVersion"`
	Name       string `yaml:"name"`
	Version    string `yaml:"version"`
}

// AppConfiguration represents the structure of the app.cfg file
type AppConfiguration struct {
	ConfigVersion string      `yaml:"app.cfg.version" json:"app.cfg.version"`
	Metadata      AppMetaData `yaml:"metadata" json:"metadata"`
}

type AppMetaData struct {
	Name        string   `yaml:"name" json:"name"`
	Icon        string   `yaml:"icon" json:"icon"`
	Description string   `yaml:"description" json:"description"`
	AppID       string   `yaml:"appid" json:"appid"`
	Title       string   `yaml:"title" json:"title"`
	Version     string   `yaml:"version" json:"version"`
	Categories  []string `yaml:"categories" json:"categories"`
	Rating      float32  `yaml:"rating" json:"rating"`
	Target      string   `yaml:"target" json:"target"`
}

func init() {
	flag.StringVar(&changedCharts, "chart-dirs", "", "comma-separated list of chart directories")
	flag.Parse()
}

func usage() {
	flag.Usage()
	os.Exit(-1)
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	dirs := strings.Split(changedCharts, ",")
	for _, dir := range dirs {
		err := validateChartFolder(dir)
		if err != nil {
			fmt.Printf("Validation failed for folder '%s': %v\n", dir, err)
			return
		}
	}

	fmt.Println("Folder validation successful!")
}

// validateChartFolder validates the specified chart folder
func validateChartFolder(folder string) error {
	// Check if the folder name is valid
	if !isValidFolderName(folder) {
		return fmt.Errorf("invalid folder name: '%s' must '^[a-z0-9]{1,30}$'", folder)
	}

	if !dirExists(folder) {
		return fmt.Errorf("folder does not exist: '%s'", folder)
	}

	// Check if Chart.yaml file exists
	chartFile := filepath.Join(folder, "Chart.yaml")
	if !fileExists(chartFile) {
		return fmt.Errorf("missing Chart.yaml in folder: '%s'", folder)
	}

	// Read and parse Chart.yaml file
	chartContent, err := os.ReadFile(chartFile)
	if err != nil {
		return fmt.Errorf("failed to read Chart.yaml in folder '%s': %v", folder, err)
	}
	var chart Chart
	if err := yaml.Unmarshal(chartContent, &chart); err != nil {
		return fmt.Errorf("failed to parse Chart.yaml in folder '%s': %v", folder, err)
	}

	// Check if Chart.yaml fields are valid
	if err := isValidChartFields(chart); err != nil {
		return err
	}

	// Check if values.yaml file exists
	valuesFile := filepath.Join(folder, "values.yaml")
	if !fileExists(valuesFile) {
		return fmt.Errorf("missing values.yaml in folder: '%s'", folder)
	}

	// Check if templates folder exists
	templatesDir := filepath.Join(folder, "templates")
	if !dirExists(templatesDir) {
		return fmt.Errorf("missing templates folder in folder: '%s'", folder)
	}

	// Check if app.cfg file exists
	appCfgFile := filepath.Join(folder, "app.cfg")
	if !fileExists(appCfgFile) {
		return fmt.Errorf("missing app.cfg in folder: '%s'", folder)
	}

	// Read and parse app.cfg file
	appCfgContent, err := os.ReadFile(appCfgFile)
	if err != nil {
		return fmt.Errorf("failed to read app.cfg in folder '%s': %v", folder, err)
	}
	var appConf AppConfiguration
	if err := yaml.Unmarshal(appCfgContent, &appConf); err != nil {
		return fmt.Errorf("failed to parse app.cfg in folder '%s': %v", folder, err)
	}

	// Check if metadata fields in app.cfg are valid
	if err := isValidMetadataFields(appConf.Metadata, chart, folder); err != nil {
		return err
	}

	// Validation passed, return nil
	return nil
}

// isValidFolderName checks if the folder name is valid using regular expression
func isValidFolderName(name string) bool {
	match, _ := regexp.MatchString("^[a-z0-9]{1,30}$", name)
	return match
}

// fileExists checks if the file exists
func fileExists(path string) bool {
	info, err := os.Stat(path)
	return (err == nil || os.IsExist(err)) && !info.IsDir()
}

// dirExists checks if the directory exists
func dirExists(path string) bool {
	info, err := os.Stat(path)
	return (err == nil || os.IsExist(err)) && info.IsDir()
}

// isValidChartFields checks if the fields in Chart.yaml are valid
func isValidChartFields(chart Chart) error {
	if chart.APIVersion == "" {
		return fmt.Errorf("apiVersion field empty in app.cfg in chart '%s'", chart)
	}

	if chart.Name == "" {
		return fmt.Errorf("name field empty in app.cfg in chart '%s'", chart)
	}

	if chart.Version == "" {
		return fmt.Errorf("version field empty in app.cfg in chart '%s'", chart)
	}

	return nil
}

// isValidMetadataFields checks if the metadata fields in app.cfg are valid
func isValidMetadataFields(metadata AppMetaData, chart Chart, folder string) error {
	if chart.Name != folder {
		return fmt.Errorf("name %s invalid in Chart.yaml in chart '%s', must same", chart.Name, folder)
	}

	if metadata.Name != folder {
		return fmt.Errorf("metadata.name %s invalid in app.cfg in chart '%s', must same", metadata.Name, folder)
	}

	if metadata.Version != chart.Version {
		return fmt.Errorf("version in app.cfg %s, version in Chart.yaml %s in chart '%s', must same", metadata.Version, chart.Version, folder)
	}

	return nil
}
