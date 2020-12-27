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

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

type Options struct {
	Silent           bool
	KeepWorkingFiles bool
	CompileOnly      bool
	PPMCKRootPath    string
	PathPppckc       string
	PathNesasm       string
	PathNsf2wav      string
	PathNesIinclude  string
	MmlFilePath      string
}

func parseOption() *Options {
	ret := &Options{}
	flag.BoolVar(&ret.Silent, "s", false, "hide output from compiles")
	flag.BoolVar(&ret.KeepWorkingFiles, "k", false, "skip cleanup ppmkc working files (define.inc, effect.h, ..)")
	flag.BoolVar(&ret.CompileOnly, "c", false, "compile mode")
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
	if opt.Silent {
		return
	}
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

func showCommandLog(opt *Options, b []byte) {
	if opt.Silent {
		return
	}

	text, err := bytesFromShiftJIS(b)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(text)
}

func envSetup(opt *Options) {
	var err error
	err = os.Setenv("PPMCK_BASEDIR", opt.PPMCKRootPath)
	if err != nil {
		log.Fatal(err)
	}
	err = os.Setenv("NES_INCLUDE", opt.PathNesIinclude)
	if err != nil {
		log.Fatal(err)
	}
}

func genNsf(opt *Options, mml, nsf, header string) {
	var ret []byte
	var err error

	ret, err = exec.Command(opt.PathPppckc, "-i", mml).CombinedOutput()
	showCommandLog(opt, ret)
	if err != nil {
		log.Fatal(err)
	}

	ret, err = exec.Command(opt.PathNesasm, "-s", "-raw", "ppmck.asm").CombinedOutput()
	showCommandLog(opt, ret)
	if err != nil {
		log.Fatal(err)
	}

	err = os.Rename("ppmck.nes", nsf)
	if err != nil {
		log.Fatal(err)
	}

	if opt.KeepWorkingFiles {
		return
	}

	for _, f := range []string{"define.inc", "effect.h", header} {
		err = os.Remove(f)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func genWave(opt *Options, src, dest string) {
	ret, err := exec.Command(opt.PathNsf2wav, src, dest).CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	showCommandLog(opt, ret)
}

func playWave(path string) {
	f, err := os.Open(path)
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
	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))
	<-done
}

func main() {
	opt := parseOption()
	showOption(opt)
	envSetup(opt)

	dir, file := filepath.Split(opt.MmlFilePath)
	mml := file
	ext := filepath.Ext(file)
	wave := strings.TrimSuffix(file, ext) + ".wav"
	nsf := strings.TrimSuffix(file, ext) + ".nsf"
	header := strings.TrimSuffix(file, ext) + ".h"

	os.Chdir(dir)

	genNsf(opt, mml, nsf, header)
	genWave(opt, nsf, wave)

	if opt.CompileOnly {
		return
	}
	playWave(wave)
}
