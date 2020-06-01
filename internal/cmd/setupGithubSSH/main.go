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
	// hostkey
	knownHosts := path.Join(sshDir, "known_hosts")
	err = ioutil.WriteFile(knownHosts, []byte("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDynhw46Nrnt7FsBg/dAGTCd1LOZTphHo0nPWidmpY/Kr/mdng/VnILpGmQa7fAlv6N9PKKm2kEUvNdnsJDLzjZch4cNLFr8Tql7k4evLBIJq7LHt6Twpc1heH6s1CGDbTZQlWDZhm/vE0jwZGH/3rjlweYQILtItMT3q6m6OQjkeLldkN5KBjHG8Fr73ucrBDc0w4ENcM7cyFYKDU8bMG2oPg86u6v0guQFgTfUydUh88ekbuIHJGvAankgrcDjnEKx2tuVBwxFyWe+Z0Q7UJW5CZVMM1ip10OQgH0CzK174reIxX2MsA0IMTWXMGsuCOJ8cBZzQqtELfrW8EunsQz root@localhost"), 0600)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
