package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

func main() {
	// Write ~/.ssh/id_rsa
	key := os.Getenv("LINODE_PRIVATE_KEY") // secret on github
	if key == "" {
		fmt.Println("LINODE_PRIVATE_KEY env not found")
		os.Exit(1)
	}
	sshDir := path.Join(os.Getenv("HOME"), ".ssh")
	err := ioutil.WriteFile(path.Join(sshDir, "id_rsa"), []byte(key), 0600)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Write ~/.ssh/known_hosts
	err = ioutil.WriteFile(
		path.Join(sshDir, "known_hosts"),
		[]byte("tidio.preferit.se ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDynhw46Nrnt7FsBg/dAGTCd1LOZTphHo0nPWidmpY/Kr/mdng/VnILpGmQa7fAlv6N9PKKm2kEUvNdnsJDLzjZch4cNLFr8Tql7k4evLBIJq7LHt6Twpc1heH6s1CGDbTZQlWDZhm/vE0jwZGH/3rjlweYQILtItMT3q6m6OQjkeLldkN5KBjHG8Fr73ucrBDc0w4ENcM7cyFYKDU8bMG2oPg86u6v0guQFgTfUydUh88ekbuIHJGvAankgrcDjnEKx2tuVBwxFyWe+Z0Q7UJW5CZVMM1ip10OQgH0CzK174reIxX2MsA0IMTWXMGsuCOJ8cBZzQqtELfrW8EunsQz\n"),
		0600,
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
