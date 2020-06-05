package perm

import "fmt"

func Example() {
	fmt.Println(NoMode)
	fmt.Println(OtherX)
	fmt.Println(UserR | UserW | GroupR | OtherR)
	// output:
	// ---------
	// --------x
	// rw-r--r--
}
