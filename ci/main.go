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
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	source := client.Container()
	source = source.From("golang:1.19")
	source = source.WithDirectory(
		"/app",
		client.Host().Directory("."),
		dagger.ContainerWithDirectoryOpts{
			Exclude: []string{"ci/"},
		},
	)

	runner := source.WithWorkdir("/app")
	out, err := runner.WithExec([]string{"go", "mod", "download"}).Stderr(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(out)

	test := runner.WithWorkdir("/app")
	out, err = test.WithExec([]string{"go", "test", "./..."}).Stderr(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(out)

	build := test.WithWorkdir("/app")
	out, err = build.WithExec([]string{"go", "build", "-v", "-o", "shopa"}).Stderr(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(out)

	dir := build.Directory("/app")
	e, err := dir.Entries(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("build dir contents:\n %s\n", e)

}
