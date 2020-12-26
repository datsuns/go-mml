package main

import (
	"flag"
	"fmt"
	"os/exec"
	"path/filepath"
)

type Options struct {
	PPMCKRootPath string
	PathPPMckc    string
}

func parseOption() *Options {
	ret := &Options{}
	flag.StringVar(&ret.PPMCKRootPath, "m", "", "path to root dir to ppmck")
	flag.Parse()

	ret.PathPPMckc = filepath.Join(ret.PPMCKRootPath, "bin", "ppmckc")
	return ret
}

func showOption(opt *Options) {
	fmt.Printf("path to ppmck  [%s]\n", opt.PPMCKRootPath)
	fmt.Printf("path to ppmckc [%s]\n", opt.PathPPMckc)
}

func main() {
	opt := parseOption()
	showOption(opt)
	//show := func(key string) {
	//	val, ok := os.LookupEnv(key)
	//	if !ok {
	//		fmt.Printf("%s not set\n", key)
	//	} else {
	//		fmt.Printf("%s=%s\n", key, val)
	//	}
	//}
	//show("PATH")

	ret, _ := exec.Command(opt.PathPPMckc, "-h").CombinedOutput()
	fmt.Println(string(ret))
}
