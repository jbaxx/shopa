package main

import (
	"context"
	"log"
	"os"

	"dagger.io/dagger"
)

func main() {

	ctx := context.Background()

	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// get reference to the local project
	src := client.Host().Directory(".")

	golang := client.Container().From("golang:latest")

	// mount cloned repository into `golang` image
	golang = golang.WithDirectory("/app", src).WithWorkdir("/app")

	// define the application build command
	path := "build/"
	golang = golang.WithExec([]string{"go", "build", "-v", "-o", path})

	// get a reference to build output directory in container
	output := golang.Directory(path)

	// write contents of container build/ directory to the host
	_, err = output.Export(ctx, path)
	if err != nil {
		log.Fatal(err)
	}

	// source := client.Container()
	// source = source.From("golang:1.19")
	// source = source.WithDirectory(
	// 	"/app",
	// 	client.Host().Directory("."),
	// 	dagger.ContainerWithDirectoryOpts{
	// 		Exclude: []string{"ci/"},
	// 	},
	// )

	// runner := source.WithWorkdir("/app")
	// out, err := runner.WithExec([]string{"go", "mod", "download"}).Stderr(ctx)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(out)

	// test := runner.WithWorkdir("/app")
	// out, err = test.WithExec([]string{"go", "test", "./..."}).Stderr(ctx)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(out)

	// build := test.WithWorkdir("/app")
	// out, err = build.WithExec([]string{"go", "build", "-v", "-o", "shopa"}).Stderr(ctx)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(out)

	// dir := build.Directory("/app")
	// e, err := dir.Entries(ctx)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("build dir contents:\n %s\n", e)

}
