package main

import (
	"fmt"
	"os"
	"strconv"
)

//
// On certain ARM64 hardware, at least on an M1 Pro performance (P) cores, this
// "pre-increment" loop is significantly slower in Go 1.17 and 1.18 than in 1.16
// This loop is also slower on the efficiency cores of an M1 Pro, but far less
// than the P cores.
//
// I used git bisect to find the change that caused this
// https://go-review.googlesource.com/c/go/+/280155
//
// $ git show 9a19481acb93114948503d935e10f6985ff15843
//
// Author: David Chase <drchase@google.com>  2020-12-30 09:05:57
// Committer: David Chase <drchase@google.com>  2021-01-12 18:40:43
// Follows: go1.16beta1
// Precedes: go1.17beta1
//
//     [dev.regabi] cmd/compile: make ordering for InvertFlags more stable
//
//     Current many architectures use a rule along the lines of
//
//     // Canonicalize the order of arguments to comparisons - helps with CSE.
//     ((CMP|CMPW) x y) && x.ID > y.ID => (InvertFlags ((CMP|CMPW) y x))
//
//     to normalize comparisons as much as possible for CSE.  Replace the
//     ID comparison with something less variable across compiler changes.
//     This helps avoid spurious failures in some of the codegen-comparison
//     tests (though the current choice of comparison is sensitive to Op
//     ordering).
//
//     Two tests changed to accommodate modified instruction choice.
//
func MakeGeneratorClosurePre() func() int {

	candidate := 1

	getprime := func() int {

		for {

			candidate++

			i := 2
			for ; i < candidate; i++ {

				if candidate%i == 0 {
					break
				}
			}

			if i == candidate {
				return candidate
			}
		}
	}

	return getprime
}

func MakeGeneratorClosurePost() func() int {

	candidate := 2

	getprime := func() int {

		for {

			i := 2
			for ; i < candidate; i++ {

				if candidate%i == 0 {
					break
				}
			}

			if i == candidate {
				candidate++
				return candidate - 1
			}

			candidate++
		}
	}

	return getprime
}

func main() {
	N, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var closure func() int

	switch os.Args[2] {
	case "pre":
		{
			closure = MakeGeneratorClosurePre()
			break
		}

	case "post":
		{
			closure = MakeGeneratorClosurePost()
			break
		}

	default:
		panic(fmt.Sprintf("unknown implementation %s", os.Args[2]))
	}

	var prime int

	for i := 0; i < N; i++ {
		prime = closure()

	}
	fmt.Printf("prime=%d\n", prime)
}
