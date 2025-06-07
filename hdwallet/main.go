package main

import (
	"git.mazdax.tech/blockchain/hdwallet/cmd"
	"path"
	"runtime"
)

func main() {
	_, caller, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}
	dir := path.Dir(caller)

	cmd.Prepare(dir)
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
