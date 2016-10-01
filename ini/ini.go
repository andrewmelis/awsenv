package ini

import (
	"bufio"
	"os"
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

func MakeINIFile(path string) (INIFile, error) {
	iniFile := INIFile{}

	f, err := os.Open(path)
	if err != nil {
		return iniFile, err
	}
	defer f.Close()

	section := INISection{}
	key := INIKey{}

	scan := bufio.NewScanner(f)
	for scan.Scan() {
		l := scan.Text()

		keySegments := strings.Split(l, "=")
		key.Name = keySegments[0]
		key.Value = keySegments[1]

		section.Keys = append(section.Keys, key)
	}

	iniFile.Sections = append(iniFile.Sections, section)

	return iniFile, err
}
