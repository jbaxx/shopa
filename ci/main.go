package main

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
)

func main() {
	ctx := context.Background()
	test(ctx)
}

func test(ctx context.Context) error {
	fmt.Println("Testing with Dagger")

	goVersions := []string{"1.20", "1.21"}

	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		return err
	}
	defer client.Close()

	// get reference to the local project
	src := client.Host().Directory(".")

	for _, version := range goVersions {
		imageTag := fmt.Sprintf("golang:%s", version)
		golang := client.Container().From(imageTag)

		// mount local project into the golang image
		golang = golang.WithDirectory("/app", src)

		// install dependencies
		runner := golang.WithWorkdir("/app")
		out, err := runner.WithExec([]string{"go", "mod", "download"}).Stderr(ctx)
		if err != nil {
			return err
		}
		fmt.Println(out)

		// run tests
		test := runner.WithWorkdir("/app")
		out, err = test.WithExec([]string{"go", "test", "./..."}).Stderr(ctx)
		if err != nil {
			return err
		}
		fmt.Println(out)

		// run vulnerability checks
		vuln := test.WithWorkdir("/app")
		out, err = vuln.WithExec([]string{
			"go",
			"install",
			"golang.org/x/vuln/cmd/govulncheck@latest",
		}).Stderr(ctx)
		if err != nil {
			return err
		}
		fmt.Println(out)

		out, err = test.WithExec([]string{"govulncheck", "./..."}).Stderr(ctx)
		if err != nil {
			return err
		}
		fmt.Println(out)

	}

	return nil
}
