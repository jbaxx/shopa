package main

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
)

func main() {
	os.Setenv("GITHUB_STEP_SUMMARY", "## This is hello! :rocket:")

	ctx := context.Background()
	if err := test(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}

func test(ctx context.Context) error {
	fmt.Println("Testing with Dagger")

	goVersions := []string{"1.20"}
	const workDir = "/app"

	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		return fmt.Errorf("dagger connect: %w", err)
	}
	defer client.Close()

	depCache := client.CacheVolume("node")

	// get reference to the local project
	src := client.Host().Directory(".")

	for _, version := range goVersions {
		imageTag := fmt.Sprintf("golang:%s", version)
		golang := client.Container().From(imageTag)

		// mount local project into the golang image
		golang = golang.WithDirectory(workDir, src).
			WithWorkdir(workDir).
			WithMountedCache("/app/node", depCache).
			WithEnvVariable("GOBIN", "$GOPATH/bin", dagger.ContainerWithEnvVariableOpts{
				Expand: true,
			})
			// WithEnvVariable("PATH", "$PATH:$HOME/go/bin", dagger.ContainerWithEnvVariableOpts{
			// 	Expand: true,
			// })

		dir := golang.Directory(workDir)
		e, err := dir.Entries(ctx)
		if err != nil {
			return fmt.Errorf("dagger entries: %w", err)
		}
		fmt.Printf("Contents of work dir %s:\n%s\n", workDir, e)

		// install dependencies
		_, err = golang.WithExec([]string{"go", "mod", "download"}).Stderr(ctx)
		if err != nil {
			return fmt.Errorf("dagger dependencies install: %w", err)
		}

		// install dependencies
		_, err = golang.WithExec([]string{"go", "mod", "tidy"}).Stdout(ctx)
		if err != nil {
			return fmt.Errorf("dagger dependencies install: %w", err)
		}

		// run tests
		_, err = golang.WithExec([]string{"go", "test", "./..."}).Stderr(ctx)
		if err != nil {
			return fmt.Errorf("dagger tests: %w", err)
		}

		// run vulnerability checks
		_, err = golang.WithExec([]string{
			"go",
			"install",
			"golang.org/x/vuln/cmd/govulncheck@latest",
		}).Stderr(ctx)
		if err != nil {
			return fmt.Errorf("dagger govulncheck install: %w", err)
		}

		dirs, err := golang.Directory("/").Entries(ctx)
		if err != nil {
			return err
		}
		fmt.Printf("Contents of / dir:\n%s\n", dirs)

		dirs, err = golang.Directory("/go").Entries(ctx)
		if err != nil {
			return err
		}
		fmt.Printf("Contents of /go dir:\n%s\n", dirs)

		dirs, err = golang.Directory("/go/bin").Entries(ctx)
		if err != nil {
			return err
		}
		fmt.Printf("Contents of /go/bin dir:\n%s\n", dirs)

		dirs, err = golang.Directory("/home").Entries(ctx)
		if err != nil {
			return err
		}
		fmt.Printf("Contents of /home dir:\n%s\n", dirs)

		_, err = golang.WithExec([]string{"go", "env"}).Stdout(ctx)
		if err != nil {
			return fmt.Errorf("dagger govulncheck install: %w", err)
		}

		path, err := golang.EnvVariable(ctx, "PATH")
		if err != nil {
			return err
		}
		fmt.Println("PATH: ", path)

		_, err = golang.WithExec([]string{"govulncheck", "./..."}).Stderr(ctx)
		if err != nil {
			return fmt.Errorf("dagger govulncheck: %w", err)
		}

	}

	return nil
}
