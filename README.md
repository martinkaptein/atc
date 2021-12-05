# ATC Air traffic controller game

This is a fork of [ndecker/atc](https://github.com/ndecker/atc), which is itself based on [Fastie ATC](https://atc.fastie.com/).

## Installation

### Instructions for MacOS and Linux

In order to use the provided binaries, find the binary of your choice inside [/binaries](/binaries).

Then download the binary, make it executable and run it:

(example)
```
wget https://raw.githubusercontent.com/martinkaptein/atc/master/binaries/atc_macOS-64bit
chmod +x atc_osx_x64
./atc_macOS-64bit
```
 
Alternatively place it in your path to run it everywhere from the terminal:
```
cp atc_macOS-64bit /usr/local/bin/atc
```

*Subsitute atc_macOS-64bit with your binary name according to your platform.*

Now you can run it by just typing `atc` in the terminal.

If any of the commands don't work, put a `sudo` before.
Keep in mind that the binary builds may be outdated, for latest features please compile.

## Building and compiling

Install golang for your platform from [the Golang website](https://golang.org/).

Clone this repo:

`git clone https://github.com/martinkaptein/atc.git`
`cd atc`

Run `go build`.

If there are errors with missing stuff, run `go get github.com/missingstuff`, then re-run `go build`.

A `atc` should appear in the working folder.

Make it executable and run it:

`chmod +x atc`
`./atc`

Alternatively place it in your $PATH.

## Windows

Download the Windows executable from `wget https://raw.githubusercontent.com/martinkaptein/atc/master/binaries/atc_windows.exe`.

Windows build instructions:

```
go build
(if errors occur)
go mod init puppy
go mod tidy
```

Tested and working on cygwin.

## Manual

The original manual (pdf) for playing can be found under `/manual/ATC_Users_Guide.pdf` [here](/manual/ATC_Users_Guide.pdf).
It is from [here](https://atc.fastie.com/instructions/how-to-play-atc.php).
