package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/rand"
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

	// // List the contents of the host directory (the local)
	// fmt.Println("### List the contents of the host directory (the local)")
	// entries, err := client.Host().Directory(".", dagger.HostDirectoryOpts{
	// 	Exclude: []string{".git*"},
	// }).Entries(ctx)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// for _, e := range entries {
	// 	fmt.Println(e)
	// }

	// // Mount the current dir (or dir specified where the "." is)
	// // in the container, in a path named "/host".
	// // Then read the contents of the mounted directory in the container.
	// fmt.Println("### Mount the current dir and list")
	// contents, err := client.Container().
	// 	From("alpine:latest").
	// 	WithDirectory("/host", client.Host().Directory(".")). // *dagger.Directory
	// 	WithExec([]string{"ls", "/host"}).
	// 	WithExec([]string{"tree", "/host"}).
	// 	Stdout(ctx)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(contents)

	// // 1. Mount the current dir (or dir specified where the "." is)
	// // in the container, in a path named "/host".
	// // 2. Write a file in the container.
	// // 3. Write that file back to the host. To write modifications back to the,
	// // host directory, you must explicitly export the directory back to the host filesystem.
	// fmt.Println("### Mount the current dir and list")
	// contenido, err := client.Container().
	// 	From("alpine:latest").
	// 	WithDirectory("/host", client.Host().Directory(".")).
	// 	WithExec([]string{"/bin/sh", "-c", `echo foo > /host/bar.txt`}).
	// 	File("/host/bar.txt").
	// 	// Directory("/host").
	// 	Export(ctx, "./bart.txt")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(contenido)

	// // add a repository to a container.
	// // get repository at specified branch
	// project := client.
	// 	Git("https://github.com/dagger/dagger").
	// 	Branch("main").
	// 	Tree()
	//
	// contents, err := client.Container().
	// 	From("alpine:latest").
	// 	WithDirectory("/src", project).
	// 	WithWorkdir("/src").
	// 	WithExec([]string{"ls", "/src"}).
	// 	Stdout(ctx)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(contents)

	// // Perform a multistage build
	// project := client.Host().Directory(".")
	// // build app
	// builder := client.Container().
	// 	From("golang:latest").
	// 	WithDirectory("/src", project).
	// 	WithWorkdir("/src").
	// 	WithEnvVariable("CGO_ENABLED", "0").
	// 	WithExec([]string{"go", "build", "-o", "myapp"})
	// // publish binary on alpine base
	// prodImage := client.Container().
	// 	From("alpine").
	// 	WithFile("/bin/myapp", builder.File("/src/myapp")).
	// 	WithEntrypoint([]string{"/bin/myapp"})
	// // addr, err := prodImage.Publish(ctx, "localhost:5000/multistage")
	// address := fmt.Sprintf("ttl.sh/hello-dagger-%.0f", math.Floor(rand.Float64()*10000000))
	// addr, err := prodImage.Publish(ctx, address)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("Address: ", address)
	// fmt.Println("addr: ", addr)

	contextDir := client.Host().Directory(".")

	address := fmt.Sprintf("ttl.sh/hello-dagger-%.0f:300s", math.Floor(rand.Float64()*10000000))
	ref, err := contextDir.
		DockerBuild().
		Publish(ctx, address)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Address: ", address)
	fmt.Printf("Publoshed image to %s\n", ref)
}
