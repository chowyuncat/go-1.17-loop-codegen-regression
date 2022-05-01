package main

import (
	"fmt"
	"os"
	"strconv"
)

//
// On certain ARM64 hardware, at least on M1 Pro performance (P) cores, this
// "pre-increment" loop is significantly slower in Go 1.17 and 1.18 than in 1.16
// This loop might also be a tiny bit slower on the M1 efficiency cores than on
// the P cores, but I have not been able to reliably measure a difference.
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
func GetNextPrimePreIncrement(candidate int) int {
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

func GetNextPrimePostIncrement(candidate int) int {

	candidate++

	for {

		i := 2
		for ; i < candidate; i++ {

			if candidate%i == 0 {
				break
			}
		}

		if i == candidate {
			return candidate
		}

		candidate++
	}
}

func main() {
	N, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}


	var getnextprime func(int) int

	switch os.Args[2] {
		case "pre": {
			getnextprime = GetNextPrimePreIncrement
			break
		}

		case "post": {
			getnextprime = GetNextPrimePostIncrement
			break
		}

		default:
			panic(fmt.Sprintf("unknown implementation %s", os.Args[2]))
	}

	prime := getnextprime(N)

	fmt.Printf("prime=%d\n", prime)
}
