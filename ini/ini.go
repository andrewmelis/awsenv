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
	matchINIProperty = regexp.MustCompile(`^\w+=\w+$`)
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

	var foundSection, foundKey bool

	scan := bufio.NewScanner(f)
	for scan.Scan() {
		l := scan.Text()

		switch {
		case l == "": // use matchstring?
			appendAndReset(&foundSection, &foundKey, &iniFile.Sections, &section, &key)
		case matchINISection.MatchString(l):
			if foundSection { // new section; no intervening newline
				appendAndReset(&foundSection, &foundKey, &iniFile.Sections, &section, &key)
			}

			section.Name = strings.Trim(l, "[]")
			foundSection = true
		case matchINIProperty.MatchString(l):
			keySegments := strings.Split(l, "=")
			key.Name = keySegments[0]
			key.Value = keySegments[1]

			foundKey = true
			section.Keys = append(section.Keys, key)
		default:
			err = fmt.Errorf("invalid INI file syntax for %s\n", l)
		}
	}

	// final section append
	appendAndReset(&foundSection, &foundKey, &iniFile.Sections, &section, &key)

	return iniFile, err
}

func appendAndReset(foundSection, foundKey *bool, sections *[]INISection, newSection *INISection, key *INIKey) {
	if *foundSection || *foundKey {
		*sections = append(*sections, *newSection) // final section append?
		// reset
		*foundSection, *foundKey = false, false
		*newSection = INISection{}
		*key = INIKey{}
	}
}

func (f INIFile) Section(name string) (INISection, error) {
	for _, section := range f.Sections {
		if section.Name == name {
			return section, nil
		}
	}
	return INISection{}, fmt.Errorf("Section %s not found!\n", name)
}
