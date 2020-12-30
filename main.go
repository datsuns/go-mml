package main

import (
	"C"
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
	"github.com/faiface/beep/mp3"
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
	PathLame         string
	PathNesIinclude  string
	MmlFilePath      string
}

func parseOption() *Options {
	ret := &Options{}
	flag.BoolVar(&ret.Silent, "s", false, "hide output from compiles")
	flag.BoolVar(&ret.KeepWorkingFiles, "k", false, "skip cleanup ppmkc working files (define.inc, effect.h, ..)")
	flag.BoolVar(&ret.CompileOnly, "c", false, "compile mode")
	flag.StringVar(&ret.PPMCKRootPath, "m", "", "path to root dir to ppmck")
	flag.StringVar(&ret.PathLame, "l", "", "path to lame command")
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
	fmt.Printf("path to lame   [%s]\n", opt.PathLame)
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
	ret := nsf2wav(src, dest)
	if ret != 0 {
		log.Fatalf("nsf2wave error [%d]\n", ret)
	}
}

func genMp3(opt *Options, src, dest string) {
	ret, err := exec.Command(opt.PathLame, src, dest).CombinedOutput()
	showCommandLog(opt, ret)
	if err != nil {
		log.Fatal(err)
	}
}

func playSound(opt *Options, path string) {
	var err error
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var format beep.Format
	var streamer beep.StreamSeekCloser
	if opt.PathLame != "" {
		streamer, format, err = mp3.Decode(f)
	} else {
		streamer, format, err = wav.Decode(f)
	}
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
	time.Sleep(time.Millisecond * 100)
}

func main() {
	opt := parseOption()
	showOption(opt)
	envSetup(opt)

	dir, file := filepath.Split(opt.MmlFilePath)
	mml := file
	ext := filepath.Ext(file)
	body := strings.TrimSuffix(file, ext)
	wavePath := body + ".wav"
	mp3Path := body + ".mp3"
	nsf := body + ".nsf"
	header := body + ".h"

	os.Chdir(dir)

	genNsf(opt, mml, nsf, header)
	genWave(opt, nsf, wavePath)
	if opt.PathLame != "" {
		genMp3(opt, wavePath, mp3Path)
		os.Remove(wavePath)
	}

	if opt.CompileOnly {
		return
	}

	if opt.PathLame != "" {
		playSound(opt, mp3Path)
	} else {
		playSound(opt, wavePath)
	}
}
