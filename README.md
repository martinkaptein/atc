# ATC Air traffic controller game

This is a fork of [ndecker/atc](https://github.com/ndecker/atc), which is itself based on [Fastie ATC](https://atc.fastie.com/).

## Installation

### Instructions for MacOS and Linux

In order to use the provided binaries, find the binary of your choice inside [/binaries](/binaries).

Then download the binary, make it executable and run it:

(example)
```
wget https://raw.githubusercontent.com/martinkaptein/atc/master/binaries/atc_osx_x64
chmod +x atc_osx_x64
./atc_osx_x64
```
 
Alternatively place it in your path to run it everywhere from the terminal:
```
cp atc_osx_x64 /usr/local/bin/atc
```

*Subsitute atc_osx_x64 with your binary name.*

Now you can run it by just typing `atc` in the terminal.

If any of the commands don't work, put a `sudo` before.

## Building

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

## Manual

The original manual (pdf) for playing can be found under `/manual/ATC_Users_Guide.pdf` [here](/manual/ATC_Users_Guide.pdf).
It is from [here](https://atc.fastie.com/instructions/how-to-play-atc.php).
