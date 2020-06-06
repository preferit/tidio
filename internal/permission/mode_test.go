package permission

import "fmt"

func Example_modes() {
	fmt.Println(NoMode)
	fmt.Println(OtherX)
	fmt.Println(UserR | UserW | GroupR | OtherR)
	// output:
	// ---------
	// --------x
	// rw-r--r--
}
