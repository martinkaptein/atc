package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var (
	DIFFICULTIES = []*Difficulty{
		&Difficulty{"Beginner", 80 * Minutes, 26},
		&Difficulty{"Easy", 60 * Minutes, 26},
		&Difficulty{"Average", 40 * Minutes, 26},
		&Difficulty{"Hard", 30 * Minutes, 26},
		&Difficulty{"Expert", 20 * Minutes, 26},
		&Difficulty{"Impossible", 16 * Minutes, 26},
	}

	BOARDS []*Board = LoadBoards("games/*.json")
	DEFAULT_BOARD *Board = BOARDS[0]
)

type Board struct {
	name string

	width  int
	height int

	entrypoints map[rune]*EntryPoint
	navaids     []Position
	routes      []Route
	nofly       []Position
}

func (b Board) Contains(p Position) bool {
	if p.x < 0 || p.y < 0 {
		return false
	}
	if p.x > b.width-1 || p.y > b.height-1 {
		return false
	}
	return true
}

func (b *Board) GetNavaid(p Position) *Position {
	for _, navaid := range b.navaids {
		if navaid == p {
			return &navaid
		}
	}
	return nil
}

func (b *Board) GetEntryPoint(p Position) *EntryPoint {
	for _, ep := range b.entrypoints {
		if ep.Position == p {
			return ep
		}
	}
	return nil
}

type EntryPoint struct {
	sign rune
	Position
	Direction
	is_airport bool
}

type Route struct {
	entry rune
	exit  rune
	Direction
	weight int
}

func (r Route) String() string {
	return fmt.Sprintf("%s-%s", string(r.entry), string(r.exit))
}

func LoadBoards(filepattern string) []*Board {
	boards := make([]*Board, 0)

	flist, err := filepath.Glob(filepattern)
	if err != nil {
		log.Fatalf("could not list files in %v\n", filepattern)
	}

	for _, f := range flist {
		boards = append(boards, ParseBoard(f))
	}

	return boards
}

func ParseBoard(filename string) *Board {
	// transfer structs for loading and converting the JSON board
	type fileRoute struct {
		Weight    int
		Entry     string
		Exit      string
		Direction string
	}

	type fileBoard struct {
		Name   string
		Board  []string
		Routes []fileRoute
	}

	buf, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("could not load board from file '%v'\n", filename)
	}

	var fBoard fileBoard
	err = json.Unmarshal(buf, &fBoard)
	if err != nil {
		log.Fatalf("failed to unmarshal board from file '%v': %v\n", filename, err)
	}

	b := Board{
		name:        fBoard.Name,
		routes:      make([]Route, 0),
		entrypoints: make(map[rune]*EntryPoint, 0),
		navaids:     make([]Position, 0),
		nofly: 		 make([]Position, 0),
	}

	b.height = len(fBoard.Board)
	if b.height <= 0 {
		log.Fatalln("board has no height")
	}

	// load board from file
	for y, row := range fBoard.Board {
		if b.width == 0 {
			b.width = len(row)
		} else if b.width != len(row) {
			log.Fatalf("inconsistent width '%v': %v, expected %v\n", row, len(row), b.width)
		}

		for x, ch := range row {
			pos := Position{x: x, y: y}

			switch ch {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				b.entrypoints[rune(ch)] = &EntryPoint{
					sign:       rune(ch),
					Position:   pos,
					is_airport: false,
				}
			case '%', '=':
				// find direction for airport
				var dir Direction
				for _, d := range DIRECTIONS {
					pos2 := pos.Move(d, 1)
					if fBoard.Board[pos2.y][pos2.x] == '+' {
						dir = d
					}
				}
				b.entrypoints[rune(ch)] = &EntryPoint{
					sign:       rune(ch),
					Position:   pos,
					Direction:  dir,
					is_airport: true,
				}
			case '+':
				// direction marker for Airport
			case '*':
				b.navaids = append(b.navaids, pos)
			case 'x':
				b.nofly = append(b.nofly, pos)
			case '.':
			default:
				panic("unknown board spec: " + string(ch))
			}
		}
	}

	// load routes
	for _, r := range fBoard.Routes {
		entry := rune(r.Entry[0])
		exit := rune(r.Exit[0])

		_, ok_entry := b.entrypoints[entry]
		_, ok_exit := b.entrypoints[exit]
		if !ok_entry || !ok_exit {
			panic("unknown entrypoint: " + string(r.Entry) + " or " + string(r.Exit))
		}

		b.routes = append(b.routes, Route{
			entry:     entry,
			exit:      exit,
			weight:    r.Weight,
			Direction: ParseDirection(r.Direction),
		})
	}

	return &b
}