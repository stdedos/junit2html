package examples

import "log"

func foo() {
	log.Println("foo")
}

func bar() { //revive:disable:unused Intentionally not covered
	log.Println("bar")
}
