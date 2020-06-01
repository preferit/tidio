package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

func main() {
	sshDir := path.Join(os.Getenv("HOME"), ".ssh")
	filename := path.Join(sshDir, "linode_rsa")
	key := os.Getenv("LINODE_PRIVATE_KEY") // secret on github
	if key == "" {
		fmt.Println("LINODE_PRIVATE_KEY env not found")
		os.Exit(1)
	}
	err := ioutil.WriteFile(filename, []byte(key), 0600)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
