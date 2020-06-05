package perm

import "fmt"

func Example() {
	fmt.Println(NoMode)
	fmt.Println(OtherExec)
	fmt.Println(UserRead | UserWrite | GroupRead | OtherRead)
	// output:
	// ---------
	// --------x
	// rw-r--r--
}
