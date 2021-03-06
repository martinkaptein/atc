package main

import (
	"fmt"
)

const (
	DIR_N   = Direction(0)
	DIR_NE  = Direction(iota)
	DIR_E   = Direction(iota)
	DIR_SE  = Direction(iota)
	DIR_S   = Direction(iota)
	DIR_SW  = Direction(iota)
	DIR_W   = Direction(iota)
	DIR_NW  = Direction(iota)
	DIR_MAX = iota
)

var DIRECTIONS = []Direction{DIR_N, DIR_NE, DIR_E, DIR_SE, DIR_S, DIR_SW, DIR_W, DIR_NW}

type Direction uint

func (d Direction) Right(n int) Direction {
	for n < 0 {
		n += DIR_MAX
	}
	return Direction((int(d) + n) % DIR_MAX)
}

func (d Direction) Left(n int) Direction {
	return d.Right(-n)
}

func (d Direction) Reverse() Direction {
	return d.Right(4)
}

func (d Direction) String() string {
	switch d {
	case DIR_N:
		return "N"
	case DIR_NE:
		return "NE"
	case DIR_E:
		return "E"
	case DIR_SE:
		return "SE"
	case DIR_S:
		return "S"
	case DIR_SW:
		return "SW"
	case DIR_W:
		return "W"
	case DIR_NW:
		return "NW"
	default:
		panic("invalid direction")
	}
}

func ParseDirection(s string) Direction {
	switch s {
	case "N":
		return DIR_N
	case "NE":
		return DIR_NE
	case "E":
		return DIR_E
	case "SE":
		return DIR_SE
	case "S":
		return DIR_S
	case "SW":
		return DIR_SW
	case "W":
		return DIR_W
	case "NW":
		return DIR_NW
	default:
		panic("unknown direction")
	}
}

type Position struct {
	x int
	y int
}

func (p Position) String() string {
	return fmt.Sprintf("%2d/%2d", p.x, p.y)
}

func (p Position) Move(d Direction, c int) Position {
	switch d {
	case DIR_N:
		return Position{p.x, p.y - c}
	case DIR_NE:
		return Position{p.x + c, p.y - c}
	case DIR_E:
		return Position{p.x + c, p.y}
	case DIR_SE:
		return Position{p.x + c, p.y + c}
	case DIR_S:
		return Position{p.x, p.y + c}
	case DIR_SW:
		return Position{p.x - c, p.y + c}
	case DIR_W:
		return Position{p.x - c, p.y}
	case DIR_NW:
		return Position{p.x - c, p.y - c}
	default:
		panic("invalid direction")
	}
}

func (p Position) Distance(p2 Position) int {
	// diagonal has unity distance => sqrt(2) = 1
	dx, dy := p2.x-p.x, p2.y-p.y
	return Max(Abs(dx), Abs(dy))
}
