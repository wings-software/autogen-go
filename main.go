package main

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/wings-software/autogen-go/builder"
	"github.com/wings-software/autogen-go/chroot"
	"github.com/wings-software/autogen-go/cloner"
)

func main() {
	var path string

	// extract the repository path
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	// if the path is a repository url,
	// clone the rempository into a temporary
	// directory, and then delete
	if isRemote(path) {
		temp, err := ioutil.TempDir("", "")
		if err != nil {
			log.Fatalln(err)
		}
		defer os.RemoveAll(temp)

		params := cloner.Params{
			Dir:        temp,
			Repo:       path,
			Username:   "", // not yet implemented
			Password:   "", // not yet implemented
			Privatekey: "", // not yet implemented
		}
		cloner := cloner.New(1, ioutil.Discard) // 1 depth, discard git clone logs
		cloner.Clone(context.Background(), params)

		// change the path to the temp directory
		path = temp
	}

	// create a chroot virtual filesystem that we
	// pass to the builder for isolation purposes.
	chroot, err := chroot.New(path)
	if err != nil {
		log.Fatalln(err)
	}

	// builds the pipeline configuration based on
	// the contents of the virtual filesystem.
	builder := builder.New("harness", "v1")
	out, err := builder.Build(chroot)
	if err != nil {
		log.Fatalln(err)
	}

	// output to console
	println(string(out))
}

// returns true if the string is a remote git repository.
func isRemote(s string) bool {
	return strings.HasPrefix(s, "git://") ||
		strings.HasPrefix(s, "http://") ||
		strings.HasPrefix(s, "https://") ||
		strings.HasPrefix(s, "git@")
}
