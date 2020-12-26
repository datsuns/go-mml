package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

type Options struct {
	PPMCKRootPath   string
	PathPPMckc      string
	PathNesIinclude string
	MmlFilePath     string
}

func parseOption() *Options {
	ret := &Options{}
	flag.StringVar(&ret.PPMCKRootPath, "m", "", "path to root dir to ppmck")
	flag.StringVar(&ret.MmlFilePath, "f", "", "path to mml file")
	flag.Parse()

	ret.PathPPMckc = filepath.Join(ret.PPMCKRootPath, "bin", "ppmckc")
	ret.PathNesIinclude = filepath.Join(ret.PPMCKRootPath, "nes_include")
	return ret
}

func showOption(opt *Options) {
	fmt.Printf("PPMCK_BASEDIR  [%s]\n", opt.PPMCKRootPath)
	fmt.Printf("NES_INCLUDE    [%s]\n", opt.PathNesIinclude)
	fmt.Printf("path to ppmckc [%s]\n", opt.PathPPMckc)
	fmt.Printf("mml file       [%s]\n", opt.MmlFilePath)
}

func transformEncoding(rawReader io.Reader, trans transform.Transformer) (string, error) {
	ret, err := ioutil.ReadAll(transform.NewReader(rawReader, trans))
	if err == nil {
		return string(ret), nil
	} else {
		return "", err
	}
}

// Convert an array of bytes (a valid ShiftJIS string) to a UTF-8 string
func bytesFromShiftJIS(b []byte) (string, error) {
	return transformEncoding(bytes.NewReader(b), japanese.ShiftJIS.NewDecoder())
}

func showCommandLog(b []byte) {
	log, err := bytesFromShiftJIS(b)
	if err != nil {
		panic(err)
	}
	fmt.Println(log)
}

func main() {
	opt := parseOption()
	showOption(opt)

	var err error
	err = os.Setenv("PPMCK_BASEDIR", opt.PPMCKRootPath)
	if err != nil {
		panic(err)
	}
	err = os.Setenv("NES_INCLUDE", opt.PathNesIinclude)
	if err != nil {
		panic(err)
	}

	ret, _ := exec.Command(opt.PathPPMckc, "-i", opt.MmlFilePath).CombinedOutput()
	showCommandLog(ret)
}
