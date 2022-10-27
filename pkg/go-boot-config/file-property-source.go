package gobootconfig

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const propertiesSource = "properties"
const yamlSource = "yaml"
const profilePropertiesSource = "profileProperties"
const profileYamlSource = "profileYaml"

const goConfigLocation = "go.config.location"
const goConfigName = "go.config.name"
const goProfiles = "go.profiles"
const goProfilesActive = "go.profiles.active"
const defaultConfigName = "application"

var defaultProfile = []string{"default"}
var defaultConfigFolder = []string{"configs", "."}
var propertyExtensions = []string{"properties"}
var yamlExtensions = []string{"yaml", "yml"}

type filePropertySource struct {
	source   string
	location string
	name     string
	value    interface{}
}

func (sps filePropertySource) getSource() string {
	return sps.source
}
func (sps filePropertySource) getLocation() string {
	return sps.location
}
func (sps filePropertySource) getName() string {
	return sps.name
}
func (sps filePropertySource) getValue() interface{} {
	return sps.value
}

func loadFromProfileFile() {
	profiles := retrieveProfiles()
	loadAllFrom(profiles, profilePropertiesSource, propertyExtensions, true, parseProperties)
	loadAllFrom(profiles, profileYamlSource, yamlExtensions, true, parseYaml)
}

func loadFromFile() {
	profiles := []string{""}
	loadAllFrom(profiles, propertiesSource, propertyExtensions, false, parseProperties)
	loadAllFrom(profiles, yamlSource, yamlExtensions, false, parseYaml)
}

func loadAllFrom(profiles []string, source string, extensions []string, forProfile bool, parser func(io.Reader, map[string]interface{})) {
	for _, path := range retrieveConfPath(profiles, extensions, forProfile) {
		if isHttp(path) {
			loadFrom(path, source, loadHttp, parser)
		} else {
			loadFrom(path, source, loadFile, parser)
		}
	}
}

func retrieveProfiles() []string {
	profiles, err := _getStringSlice(goProfiles)
	if err != nil {
		profiles = make([]string, len(defaultProfile))
		copy(profiles, defaultProfile)
	}
	activeProfiles, err := _getStringSlice(goProfilesActive)
	if err == nil {
		profiles = filerActiveProfiles(activeProfiles, profiles)
	}
	for i, profile := range profiles {
		profiles[i] = fmt.Sprintf("-%s", profile)
	}
	return profiles
}

func filerActiveProfiles(activeProfiles []string, profiles []string) []string {
	tmpProfiles := make([]string, 0)
	for _, activeProfile := range activeProfiles {
		for _, profile := range profiles {
			if profile == activeProfile {
				tmpProfiles = append(tmpProfiles, profile)
				continue
			}
		}
	}
	return tmpProfiles
}

func loadFrom(path, source string, loader func(string, func(io.Reader, map[string]interface{})) (map[string]interface{}, error), parser func(io.Reader, map[string]interface{})) {
	configMap, err := loader(path, parser)
	if err != nil {
		log.Fatalf("Failed loading %s: %s", path, err)
		return
	}
	for k, v := range configMap {
		_addPropertySource(k, filePropertySource{source, path, k, v})
	}
}

func loadFile(filePath string, read func(io.Reader, map[string]interface{})) (map[string]interface{}, error) {
	configs := make(map[string]interface{})
	if ok, _ := pathExists(filePath, false); ok {
		file, err := os.Open(filePath)
		if err != nil {
			log.Fatalf("Failed opening file %s: %s", filePath, err)
			return configs, err
		}
		defer file.Close()
		read(file, configs)
		return configs, nil
	}
	return configs, errors.New("path not found")
}

func loadHttp(httpPath string, read func(io.Reader, map[string]interface{})) (map[string]interface{}, error) {
	configs := make(map[string]interface{})
	resp, err := http.Get(httpPath)
	if err != nil {
		log.Fatalf("Failed opening http %s: %s", httpPath, err)
		return configs, err
	}
	defer resp.Body.Close()
	read(resp.Body, configs)
	return configs, nil
}

func retrieveConfPath(profiles []string, extensions []string, forProfile bool) []string {
	paths := make([]string, 0)
	locations, err := _getStringSlice(goConfigLocation)
	if err != nil {
		locations = make([]string, len(defaultConfigFolder))
		copy(locations, defaultConfigFolder)
	}
	appName, err := _getString(goConfigName)
	if err != nil {
		appName = defaultConfigName
	}
	for _, location := range locations {
		isDir, e := pathExists(location, true)
		if e != nil {
			if !forProfile && isHttp(location) {
				for _, extension := range extensions {
					if fmt.Sprintf(".%s", extension) == filepath.Ext(location) {
						paths = append(paths, location)
						continue
					}
				}
			}
			continue
		}
		if !isDir && !forProfile {
			for _, extension := range extensions {
				if fmt.Sprintf(".%s", extension) == filepath.Ext(location) {
					paths = append(paths, location)
					continue
				}
			}
		} else {
			for _, profile := range profiles {
				for _, extension := range extensions {
					filePath := filepath.Join(location, fmt.Sprintf("%s%s.%s", appName, profile, extension))
					isFile, e := pathExists(filePath, false)
					if e != nil {
						continue
					}
					if isFile {
						paths = append(paths, filePath)
					}
				}
			}
		}
	}
	return paths
}

func isHttp(location string) bool {
	return strings.HasPrefix(strings.ToLower(location), "http")
}

func pathExists(filePath string, isDir bool) (bool, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return false, err
	}
	return info.IsDir() == isDir, nil
}
