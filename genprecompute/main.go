// genprecompute computes base-point doubles for the secp256k1 generator point,
// and writes them as generated code to a given output file.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kklash/ekliptic"
)

var (
	outputFilePath string
)

func run() error {
	flag.StringVar(&outputFilePath, "o", "", "output precomputed doubles code to this file")
	flag.Parse()

	doubles := ekliptic.ComputePointDoubles(
		ekliptic.Secp256k1_GeneratorX,
		ekliptic.Secp256k1_GeneratorY,
	)

	code := `package ekliptic

// DO NOT EDIT!
// This file was automatically generated.

func init() {
	basePointPrecomputations = PrecomputedDoubles{
`

	for _, double := range doubles {
		code += "\t\t{\n"
		code += fmt.Sprintf("\t\t\thexint(\"%x\"),\n", double[0])
		code += fmt.Sprintf("\t\t\thexint(\"%x\"),\n", double[1])
		code += "\t\t},\n"
	}
	code += "\t}\n"
	code += "}\n"

	return os.WriteFile(outputFilePath, []byte(code), 0644)
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
}
