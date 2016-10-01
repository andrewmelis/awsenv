package main

import (
	"flag"
	"fmt"
)

func main() {
	profile := flag.String("profile", "default", "aws profile to use")
	credsLocation := flag.String("location", "~/.aws/credentials", "location of aws credentials file")
	flag.Parse()

	// awsCredentials := MakeINIFile(*credsLocation)

	// fmt.Printf("opts: %s profile at %s\n%#v\n", *profile, *credsLocation, awsCredentials)
	fmt.Printf("opts: %s profile at %s\n", *profile, *credsLocation)
}
