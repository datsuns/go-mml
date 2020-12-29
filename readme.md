# play MML(Music Macro Language)

## How it works

1. use `ppmck` to compile mml file to nsf.
1. convert nsf file to wave by `Game Music Emu` library.
   * use `Game Music Emu` by cgo
1. play wave file.

## setup

1. install mml compiler
   * visit to http://ppmck.web.fc2.com/ppmck.html
1. install dependent package
   * `go get -u github.com/faiface/beep`
   * `go get -u golang.org/x/text`
1. then, download this package, and `make`
   * `git clone https://github.com/datsuns/go-mml`
   * `cd go-mml && make lib && make`

## usage

`go-mml [-s][-k][-c][-l <lame-command>] -m <ppmck-basedir> -f <mml-filepath>`

* -m: root directory to location of ppmck.
   * e.g. `c:\Users\user\work\ppmck\mck`
* -f: path to mml file
* -s: silent mode. hide outputs from commands. default: **FALSE**
* -k: keep working files. skip cleanup ppmck working files. default: **FALSE**
* -c: compile only. skip playing wave file. default: **FALSE**
* -l: path to lame command. play by mp3 if this option specified.

# appendix

* https://github.com/datsuns/vim-mml-play
