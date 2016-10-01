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

func makeTestFile(names, values []string) *os.File {
	if len(names) != len(values) {
		log.Fatal()
	}

	var content string
	for i, _ := range names {
		if i == 0 {
			content = fmt.Sprintf("%s=%s\n", names[i], values[i])
		} else {
			content = fmt.Sprintf("%s%s=%s\n", content, names[i], values[i])
		}
	}

	testFile, err := ioutil.TempFile("", "ini")
	if err != nil {
		log.Fatal()
	}
	fmt.Printf("content:\n%+v\n", content)
	testFile.Write([]byte(content))
	return testFile
}

func TestMakeINIFileNoSectionSingleKey(t *testing.T) {
	names := []string{"testName"}
	values := []string{"testValue"}
	testFile := makeTestFile(names, values)
	defer os.Remove(testFile.Name())

	ini, err := MakeINIFile(testFile.Name())
	if err != nil {
		t.Fail()
	}

	if ini.Sections[0].Name != "" {
		t.Fail()
	}

	if len(ini.Sections[0].Keys) != len(names) {
		t.Fail()
	}

	key := ini.Sections[0].Keys[0]
	if key.Name != "testName" && key.Value != "testValue" {
		t.Fail()
	}
}

func TestMakeINIFileNoSectionMultipleKeys(t *testing.T) {
	testNames := []string{"testName", "name2"}
	testValues := []string{"testValue", "value2"}
	testFile := makeTestFile(testNames, testValues)
	defer os.Remove(testFile.Name())

	ini, err := MakeINIFile(testFile.Name())
	if err != nil {
		t.Fail()
	}
	fmt.Printf("%+v\n", ini)

	if ini.Sections[0].Name != "" {
		t.Fail()
	}

	if len(ini.Sections[0].Keys) != len(testNames) {
		t.Fail()
	}

	for i, key := range ini.Sections[0].Keys {
		fmt.Printf("key: %+v\ntestNames[%d]: %s, testValues[%d]: %s\n", key, i, testNames[i], i, testValues[i])
		if key.Name != testNames[i] && key.Value != testValues[i] {
			t.Fail()
		}
	}
}
