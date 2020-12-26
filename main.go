package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

type Options struct {
	PPMCKRootPath   string
	PathPppckc      string
	PathNesasm      string
	PathNsf2wav     string
	PathNesIinclude string
	MmlFilePath     string
}

func parseOption() *Options {
	ret := &Options{}
	flag.StringVar(&ret.PPMCKRootPath, "m", "", "path to root dir to ppmck")
	flag.StringVar(&ret.PathNsf2wav, "n", "", "path to nsf2wav command")
	flag.StringVar(&ret.MmlFilePath, "f", "", "path to mml file")
	flag.Parse()

	ret.PathPppckc = filepath.Join(ret.PPMCKRootPath, "bin", "ppmckc")
	ret.PathNesasm = filepath.Join(ret.PPMCKRootPath, "bin", "nesasm")
	ret.PathNesIinclude = filepath.Join(ret.PPMCKRootPath, "nes_include")
	return ret
}

func showOption(opt *Options) {
	fmt.Printf("PPMCK_BASEDIR  [%s]\n", opt.PPMCKRootPath)
	fmt.Printf("NES_INCLUDE    [%s]\n", opt.PathNesIinclude)
	fmt.Printf("path to ppmckc [%s]\n", opt.PathPppckc)
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

	dir, file := filepath.Split(opt.MmlFilePath)
	ext := filepath.Ext(file)
	nsf := strings.TrimSuffix(file, ext) + ".nsf"
	wave := strings.TrimSuffix(file, ext) + ".wav"
	header := strings.TrimSuffix(file, ext) + ".h"

	os.Chdir(dir)

	var ret []byte
	ret, _ = exec.Command(opt.PathPppckc, "-i", opt.MmlFilePath).CombinedOutput()
	showCommandLog(ret)

	ret, _ = exec.Command(opt.PathNesasm, "-s", "-raw", "ppmck.asm").CombinedOutput()
	showCommandLog(ret)

	err = os.Rename("ppmck.nes", nsf)
	if err != nil {
		panic(err)
	}

	for _, f := range []string{"define.inc", "effect.h", header} {
		err = os.Remove(f)
		if err != nil {
			panic(err)
		}
	}

	ret, _ = exec.Command(opt.PathNsf2wav, nsf, wave).CombinedOutput()
	showCommandLog(ret)

	f, err := os.Open(wave)
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := wav.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()
	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		log.Fatal(err)
	}
	speaker.Play(streamer)
	select {}
}
