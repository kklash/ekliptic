// genprecompute precomputes a base-point table for the secp256k1 generator point,
// and writes it as generated code to a given output file.
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
	flag.StringVar(&outputFilePath, "o", "", "output precomputed table code to this file")
	flag.Parse()

	table := ekliptic.NewPrecomputedTable(
		ekliptic.Secp256k1_GeneratorX,
		ekliptic.Secp256k1_GeneratorY,
	)

	code := `package ekliptic

import "math/big"

// DO NOT EDIT!
// This file was automatically generated.

func init() {
	basePointPrecomputations = PrecomputedTable{
`

	for _, row := range table {
		code += "\t\t{\n"
		for _, point := range row {
			code += "\t\t\t{\n"
			code += fmt.Sprintf("\t\t\t\tnew(big.Int).SetBytes(%#v),\n", point[0].Bytes())
			code += fmt.Sprintf("\t\t\t\tnew(big.Int).SetBytes(%#v),\n", point[1].Bytes())
			code += "\t\t\t},\n"
		}
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
