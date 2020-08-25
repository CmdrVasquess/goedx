package l10n

import (
	"fmt"
)

func ExampleNormKey() {
	fmt.Println(normKey("foo"))
	fmt.Println(normKey("$bar"))
	fmt.Println(normKey("baz;"))
	fmt.Println(normKey("$quux;"))
	// Output:
	// foo
	// bar
	// baz
	// quux
}
