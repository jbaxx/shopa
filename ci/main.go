package main

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
)

func main() {
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

	// get reference to the local project
	src := client.Host().Directory(".")

	for _, version := range goVersions {
		imageTag := fmt.Sprintf("golang:%s", version)
		golang := client.Container().From(imageTag)

		// mount local project into the golang image
		golang = golang.WithDirectory(workDir, src)

		dir := golang.Directory(workDir)
		e, err := dir.Entries(ctx)
		if err != nil {
			return fmt.Errorf("dagger entries: %w", err)
		}
		fmt.Printf("Contents of work dir %s:\n%s\n", workDir, e)

		// install dependencies
		runner := golang.WithWorkdir(workDir)
		out, err := runner.WithExec([]string{"go", "mod", "download"}).Stderr(ctx)
		if err != nil {
			return fmt.Errorf("dagger dependencies install: %w", err)
		}
		fmt.Println(out)

		// run tests
		test := runner.WithWorkdir(workDir)
		out, err = test.WithExec([]string{"go", "test", "./..."}).Stderr(ctx)
		if err != nil {
			return fmt.Errorf("dagger tests: %w", err)
		}
		fmt.Println(out)

		// run vulnerability checks
		vuln := test.WithWorkdir(workDir)
		out, err = vuln.WithExec([]string{
			"go",
			"install",
			"golang.org/x/vuln/cmd/govulncheck@latest",
		}).Stderr(ctx)
		if err != nil {
			return fmt.Errorf("dagger govulncheck install: %w", err)
		}
		fmt.Println(out)

		out, err = test.WithExec([]string{"govulncheck", "./..."}).Stderr(ctx)
		if err != nil {
			return fmt.Errorf("dagger govulncheck: %w", err)
		}
		fmt.Println(out)

	}

	return nil
}
