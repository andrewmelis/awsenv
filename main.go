package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/andrewmelis/awsenv/ini"
)

func main() {
	profile := flag.String("profile", "default", "aws profile to use")
	credsLocation := flag.String("location", "~/.aws/credentials", "location of aws credentials file")
	flag.Parse()

	credentialsFile, err := ini.MakeINIFile(*credsLocation)
	if err != nil {
		log.Fatal("error retrieving credentials: %s\n", err)
	}

	targetProfile, err := credentialsFile.Section(*profile)
	if err != nil {
		log.Fatal("error retrieving credentials: %s\n", err)
	}

	exportCredentials(AwsProfile(targetProfile))
	// fmt.Printf("opts: %s profile at %s\n%+v\n", *profile, *credsLocation, credentialsFile)
}

type AwsProfile ini.INISection

type ShellConfig struct {
	Credentials map[string]string
}

func generateExportCredentials(profile AwsProfile) map[string]string {
	cfg := make(map[string]string)

	fmt.Printf("profile: %+v\n", profile)

	for _, key := range profile.Keys {
		upperName := strings.ToUpper(key.Name)
		cfg[upperName] = key.Value
	}
	return cfg
}
