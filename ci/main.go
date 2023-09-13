package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"dagger.io/dagger"
)

func main() {

	ctx := context.Background()

	err := build(ctx)
	if err != nil {
		log.Fatal(err)
	}

}

func build(ctx context.Context) error {
	fmt.Println("Building with Dagger")

	// define build matrix
	oses := []string{"linux", "darwin"}
	arches := []string{"amd64", "arm64"}

	// initialize Dagger client
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		return err
	}
	defer client.Close()

	// get reference to the local project
	src := client.Host().Directory(".")

	// create empty directory to put the build outputs
	outputs := client.Directory()

	golang := client.Container().From("golang:latest")

	// mount cloned repository into `golang` image
	golang = golang.WithDirectory("/app", src).WithWorkdir("/app")

	for _, goos := range oses {
		for _, goarch := range arches {
			// create a directory for each os and arch
			path := fmt.Sprintf("build/%s/%s/", goos, goarch)

			// set GOARCH and GOOS in the build environment
			build := golang.WithEnvVariable("GOOS", goos)
			build = build.WithEnvVariable("GOARCH", goarch)

			// build application
			build = build.WithExec([]string{"go", "build", "-v", "-o", path})

			// get reference to build output directory in container
			// WithDirectory write a directory at the given path,
			// with the contents of the passed dir.
			// Directory retrieves a directory.
			// This is creating new paths (to build the tree),
			// with the contents of the current build path.
			outputs = outputs.WithDirectory(path, build.Directory(path))
		}
	}

	// write build artifacts to host
	_, err = outputs.Export(ctx, ".")
	if err != nil {
		return err
	}

	return nil
}
