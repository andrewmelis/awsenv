package ini

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestMakeINIFileNoSuchFile(t *testing.T) {
	_, err := MakeINIFile("fake")
	if err == nil {
		t.Fail()
	}
}

func TestMakeINIFileInvalidContents(t *testing.T) {
	badTestFile := makeInvalidTestFile()
	_, err := MakeINIFile(badTestFile.Name())
	if err == nil {
		t.Fail()
	}
}

func TestMakeINIFileNoSectionMultipleKeys(t *testing.T) {
	testNames := []string{"name", "name_2", "name3"}
	testValues := []string{"value", "value2", "value_3"}
	testFile := makeValidTestFile(nil, testNames, testValues)
	defer os.Remove(testFile.Name())

	ini, err := MakeINIFile(testFile.Name())
	if err != nil {
		t.Fail()
	}

	if ini.Sections[0].Name != "" {
		t.Fail()
	}

	if len(ini.Sections[0].Keys) != len(testNames) {
		t.Fail()
	}

	for i, key := range ini.Sections[0].Keys {
		if key.Name != testNames[i] && key.Value != testValues[i] {
			t.Fail()
		}
	}
}

func TestMakeINIFileSectionsMultipleKeysLineBreaks(t *testing.T) {
	testSections := []string{"section1", "section_2", "section3"}
	testNames := []string{"name", "name2", "name_3"}
	testValues := []string{"value_", "value2", "value3"}
	testFile := makeValidTestFile(testSections, testNames, testValues)
	defer os.Remove(testFile.Name())

	ini, err := MakeINIFile(testFile.Name())
	if err != nil {
		t.Errorf("error: %s", err)
	}

	for i, section := range ini.Sections {
		if section.Name != testSections[i] {
			t.Errorf("actual section name %s does not match expected %s", section.Name, testSections[i])
		}

		for j, key := range section.Keys {
			if key.Name != testNames[j] {
				t.Errorf("actual key name %s does not match expected %s", key.Name, testNames[j])
			}

			if key.Value != testValues[j] {
				t.Errorf("actual key value %s does not match expected %s", key.Value, testValues[j])
			}
		}
	}
}

func TestGetSectionValidName(t *testing.T) {
	testSections := []string{"section1", "section_2", "section3"}
	testNames := []string{"name", "name2", "name_3"}
	testValues := []string{"value_", "value2", "value3"}
	testFile := makeValidTestFile(testSections, testNames, testValues)
	defer os.Remove(testFile.Name())

	ini, err := MakeINIFile(testFile.Name())

	section, err := ini.Section("section_2")
	if err != nil {
		t.Errorf("Error retrieving section: %s\n", err)
	}

	for i, key := range section.Keys {
		if key.Name != testNames[i] {
			t.Errorf("actual key name %s does not match expected %s", key.Name, testNames[i])
		}

		if key.Value != testValues[i] {
			t.Errorf("actual key value %s does not match expected %s", key.Value, testValues[i])
		}
	}
}

func TestGetSectionInvalidName(t *testing.T) {
	testSections := []string{"section1", "section_2", "section3"}
	testNames := []string{"name", "name2", "name_3"}
	testValues := []string{"value_", "value2", "value3"}
	testFile := makeValidTestFile(testSections, testNames, testValues)
	defer os.Remove(testFile.Name())

	ini, err := MakeINIFile(testFile.Name())

	section, err := ini.Section("fake_section2")
	if err == nil {
		t.Errorf("returned incorrect section: %s\n", section)
	}
}

/*
example output:
[section1]
name=value
name2=value2
name3=value3

[section2]
name=value
name2=value2
name3=value3 <-- intentionally lacking newline
[section3]
name=value
name2=value2
name3=value3
*/
func makeValidTestFile(sections, names, values []string) *os.File {
	if len(names) != len(values) {
		log.Fatal()
	}

	if sections == nil {
		sections = []string{}
	}

	var content string

	if len(sections) < 1 {
		for i := range names {
			if i == 0 {
				content = fmt.Sprintf("%s=%s\n", names[i], values[i])
			} else {
				content = fmt.Sprintf("%s%s=%s\n", content, names[i], values[i])
			}
		}
	} else {
		for i, section := range sections {
			if i == 0 {
				content = fmt.Sprintf("[%s]\n", section)
			} else {
				content = fmt.Sprintf("%s[%s]\n", content, section)
			}

			for j := range names {
				content = fmt.Sprintf("%s%s=%s\n", content, names[j], values[j])
			}

			// no linebreak after second section for arbitrary formatting test
			if i == 1 {
				content = fmt.Sprintf("%s", content)
			} else {
				content = fmt.Sprintf("%s\n", content)
			}
		}
	}
	return makeTestFile(content)
}

func makeInvalidTestFile() *os.File {
	return makeTestFile("badstuff]")
}

func makeTestFile(content string) *os.File {
	testFile, err := ioutil.TempFile("", "ini")
	if err != nil {
		log.Fatal()
	}
	testFile.Write([]byte(content))
	return testFile
}
