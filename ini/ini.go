package ini

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type INIFile struct {
	Sections []INISection
}

type INISection struct {
	Name string
	Keys []INIKey
}

type INIKey struct {
	Name  string
	Value string // interface{}?
}

var (
	matchINISection  = regexp.MustCompile(`^\[\w+\]$`)
	matchINIProperty = regexp.MustCompile(`^[0-9A-Za-z]+=[0-9A-Za-z]+$`)
)

func MakeINIFile(path string) (INIFile, error) {
	iniFile := INIFile{}

	f, err := os.Open(path)
	if err != nil {
		return iniFile, err
	}
	defer f.Close()

	section := INISection{}
	key := INIKey{}

	var pendingSection, pendingKey bool

	scan := bufio.NewScanner(f)
	for scan.Scan() {
		l := scan.Text()

		switch {
		case l == "": // use matchstring?
			if pendingSection || pendingKey {
				iniFile.Sections = append(iniFile.Sections, section) // final section append?
				// reset
				pendingSection, pendingKey = false, false
				section = INISection{}
				key = INIKey{}
			}
		case matchINISection.MatchString(l):
			section.Name = strings.Trim(l, "[]")
			pendingSection = true
		case matchINIProperty.MatchString(l):
			keySegments := strings.Split(l, "=")
			key.Name = keySegments[0]
			key.Value = keySegments[1]

			pendingKey = true
			section.Keys = append(section.Keys, key)
		default:
			err = fmt.Errorf("invalid INI file syntax for %s\n", l)
		}
	}

	// final section append
	if pendingSection || pendingKey {
		iniFile.Sections = append(iniFile.Sections, section) // final section append?
		pendingSection, pendingKey = false, false
	}

	return iniFile, err
}
