package config

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
)

// ConfigMap is shorthand for the type used as a config struct.
type ConfigMap map[string]map[string][]string

var (
	configSection = regexp.MustCompile("^\\s*\\[\\s*(\\w+)\\s*\\]\\s*$")
	configLine    = regexp.MustCompile("^\\s*(\\w+)\\s*=\\s*(.*)\\s*$")
	commentLine   = regexp.MustCompile("^#.*$")
	blankLine     = regexp.MustCompile("^\\s*$")
)

var DefaultSection = "default"

// ParseFile takes the filename as a string and returns a ConfigMap.
func ParseFile(fileName string) (cfg ConfigMap, err error) {
	var file *os.File

	cfg = make(ConfigMap, 0)
	file, err = os.Open(fileName)
	if err != nil {
		return
	}
	defer file.Close()
	buf := bufio.NewReader(file)

	var (
		line           string
		longLine       bool
		currentSection string
		lineBytes      []byte
		isPrefix       bool
	)

	for {
		err = nil
		lineBytes, isPrefix, err = buf.ReadLine()
		if io.EOF == err {
			err = nil
			break
		}

		if err != nil {
			break
		}

		if isPrefix {
			line += string(lineBytes)
			longLine = true
			continue
		}

		if longLine {
			line += string(lineBytes)
			longLine = false
		} else {
			line = string(lineBytes)
		}

		if commentLine.MatchString(line) {
			continue
		}

		if blankLine.MatchString(line) {
			continue
		}

		if configSection.MatchString(line) {
			section := configSection.ReplaceAllString(line, "$1")
			if section == "" {
				err = fmt.Errorf("invalid structure in file")
				break
			}

			if !cfg.SectionInConfig(section) {
				cfg[section] = make(map[string][]string, 0)
			}
			currentSection = section
			continue
		}

		if configLine.MatchString(line) {
			if currentSection == "" {
				currentSection = DefaultSection
			}
			key := configLine.ReplaceAllString(line, "$1")
			val := configLine.ReplaceAllString(line, "$2")
			if key == "" {
				continue
			}
			if _, ok := cfg[currentSection][key]; !ok {
				strs := make([]string, 0)
				strs = append(strs, val)
				cfg[currentSection][key] = strs
				continue
			}
			cfg[currentSection][key] = append(cfg[currentSection][key], val)
			continue
		}

		err = fmt.Errorf("invalid config file")
		break
	}

	return
}

// SectionInConfig determines whether a section is in the configuration.
func (c *ConfigMap) SectionInConfig(section string) bool {
	_, ok := (*c)[section]
	return ok
}

// ListSections returns the list of sections in the config map.
func (c *ConfigMap) ListSections() (sections []string) {
	for section, _ := range *c {
		sections = append(sections, section)
	}
	return
}
