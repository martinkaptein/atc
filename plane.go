package main

import (
	"fmt"
	"math/rand"
	"sort"
)

type PlaneState int

const (
	StatePending  = 0
	StateIncoming = iota
	StateWaiting  = iota
	StateRolling  = iota
	StateFlying   = iota
	StateAproach  = iota

	StateLanded   = iota
	StateDeparted = iota
)

const SAFE_DISTANCE = 3
const FUEL_INDICATOR = 10 * Minutes
const REUSE_ENTRYPOINT_TIME = 3 * Minutes

type Plane struct {
	callsign rune
	typ      *PlaneType

	entry *EntryPoint
	exit  *EntryPoint

	start      Ticks
	state      PlaneState
	wait_ticks Ticks

	fuel_left Ticks

	Position
	is_hoovering bool

	Direction
	want_turn int

	height         int
	want_height    int
	last_height    int
	initial_height int

	hold_at_navaid   bool
	is_holding       bool
	clear_to_aproach rune
}

func (p *Plane) Tick(game *GameState) (er *EndReason) {
	for n := 1; n <= p.typ.moves_per_tick; n += 1 {
		er := p.DoTick(game)
		if er != nil {
			return er
		}
	}
	return nil
}

func (p *Plane) DoTick(game *GameState) *EndReason {
	if p.IsConsumingFuel() {
		p.fuel_left -= 1
		if p.fuel_left == 0 {
			return &EndReason{
				message: "Fuel exhausted",
				planes:  []*Plane{p},
			}
		}
	}

	if p.wait_ticks > 0 {
		p.wait_ticks -= 1
		return nil
	}

	p.last_height = p.height

	switch p.state {
	case StatePending:
		// Wait until visible
		if (game.clock - p.start) <= p.typ.ticks_pending {
			if p.entry.is_airport {
				p.state = StateWaiting
				p.height = 0
			} else {
				p.state = StateIncoming
				p.wait_ticks = p.typ.ticks_pending
				p.height = p.initial_height
			}

			p.want_height = p.height
			p.Position = p.entry.Position
		}
	case StateIncoming: // after wait ticks
		p.state = StateFlying

	case StateWaiting: // wait for DoHeight
	case StateRolling: // after wait ticks
		er := p.UpdatePosition(game)
		if er != nil {
			return er
		}
		p.ApplyWants()

		p.state = StateFlying
	case StateFlying, StateAproach:
		if p.is_holding {
			p.Direction = p.Direction.Left(1)
		}

		er := p.UpdatePosition(game)
		if er != nil {
			return er
		}
		p.ApplyWants()

		navaid := game.board.GetNavaid(p.Position)
		if navaid != nil {
			if p.hold_at_navaid {
				p.is_holding = true
			}

			if ep, ok := game.board.entrypoints[p.clear_to_aproach]; ok {
				// always use the direction of the airport.
				p.Direction = ep.Direction
			}
		}

		p.wait_ticks = p.typ.ticks_per_move - 1
	case StateDeparted, StateLanded:
	default:
		panic("unhandled case")
	}
	return nil
}

func (p *Plane) Collides(p2 *Plane) bool {
	height_match := false

	if p.height == p2.height {
		// same height
		height_match = true
	} else if p.height == p2.last_height && p.last_height == p2.height {
		// crossover
		height_match = true
	}

	if !height_match {
		return false
	}

	distance := p.Position.Distance(p2.Position)
	return distance < SAFE_DISTANCE
}

func (p *Plane) ApplyWants() {
	if p.want_turn > 0 {
		p.Direction = p.Direction.Right(1)
		p.want_turn -= 1
	}
	if p.want_turn < 0 {
		p.Direction = p.Direction.Left(1)
		p.want_turn += 1
	}

	if p.want_height > p.height {
		p.height += 1
	}
	if p.want_height < p.height {
		p.height -= 1
	}
}

func (p *Plane) UpdatePosition(game *GameState) *EndReason {
	next_pos := p.Position
	if !p.is_hoovering {
		next_pos = p.Position.Move(p.Direction, 1)
	}

	if !game.board.Contains(next_pos) {
		// left the playing field
		if p.Position == p.exit.Position && p.height == p.typ.exit_height {
			p.state = StateDeparted
			return nil
		} else {
			return &EndReason{
				message: "Boundary Error",
				planes:  []*Plane{p},
			}
		}
	}

	if !p.typ.can_enter_nofly {
		for _, nf := range game.board.nofly {
			if nf == next_pos {
				return &EndReason{
					message: "Entering NoFly Area",
					planes:  []*Plane{p},
				}
			}
		}
	}

	if p.state == StateAproach {
		ap := game.board.GetEntryPoint(next_pos)
		if ap != nil {
			if ap == p.exit && p.height == 0 {
				p.state = StateLanded
				return nil
			}

			// call off landing
			if !p.is_hoovering {
				// if hoovering over the airport do not reset to flying
				p.state = StateFlying
				p.height = 1
			}
		}
	}
	p.Position = next_pos
	return nil
}

func (p *Plane) DoTurn(c int) bool {
	if c < -4 || c > 4 {
		return false
	}

	if p.typ.immediate_turn {
		p.Direction = p.Direction.Right(c)
		p.want_turn = 0
	} else {
		p.want_turn = c
	}
	p.is_holding = false
	p.hold_at_navaid = false
	p.clear_to_aproach = 0
	return true
}

func (p *Plane) DoHeight(h int) bool {
	if h > 5 || h < 0 {
		return false
	}

	if h == 0 {
		// aproach
		if p.state != StateFlying {
			return false
		}
		p.state = StateAproach
		p.want_height = h
		return true
	}

	p.want_height = h

	if p.state == StateWaiting {
		p.state = StateRolling
		p.wait_ticks = p.typ.ticks_rolling
	}
	return true
}

func (p *Plane) DoHold() bool {
	p.hold_at_navaid = true
	p.clear_to_aproach = 0
	return true
}

func (p *Plane) DoKeep() bool {
	if !p.typ.can_hoover {
		return false
	}
	p.is_hoovering = !p.is_hoovering
	return true
}

func (p *Plane) TurnAtNavaid(navaid rune) bool {
	p.clear_to_aproach = navaid
	p.hold_at_navaid = false
	return true
}

func (p Plane) AcceptsCommands() bool {
	return p.state == StateWaiting || p.state == StateRolling || p.state == StateFlying
}
func (p Plane) IsVisible() bool {
	return p.state == StateIncoming || p.state == StateWaiting
}
func (p Plane) IsActive() bool {
	return p.state == StateIncoming ||
		p.state == StateWaiting ||
		p.state == StateRolling ||
		p.state == StateFlying ||
		p.state == StateAproach
}
func (p Plane) IsFlying() bool {
	return p.state == StateFlying || p.state == StateAproach
}
func (p Plane) IsConsumingFuel() bool {
	return p.state == StateWaiting ||
		p.state == StateRolling ||
		p.state == StateFlying ||
		p.state == StateAproach
}

func (p Plane) IsDone() bool {
	return p.state == StateLanded ||
		p.state == StateDeparted
}

func (p Plane) String() string {
	return fmt.Sprintf("%s %s",
		p.start, p.State())
}

func (p Plane) Flightplan() string {
	if p.initial_height <= 9 {
		return fmt.Sprintf("%c%d%c %c-%c",
			p.callsign, p.initial_height, p.typ.mark, p.entry.sign, p.exit.sign)
	} else {
		return fmt.Sprintf("%cX%c %c-%c",
			p.callsign, p.typ.mark, p.entry.sign, p.exit.sign)
	}
}

func (p Plane) Marker() string {
	if p.height <= 9 {
		return fmt.Sprintf("%c%d", p.callsign, p.height)
	} else {
		return fmt.Sprintf("%cX", p.callsign)
	}
}

func (p Plane) State() string {
	res := fmt.Sprintf("%s%c %c-%c %-2s",
		p.Marker(), p.typ.mark, p.entry.sign, p.exit.sign, p.Direction)

	if p.is_hoovering {
		res += " H "
	}

	if p.fuel_left >= FUEL_INDICATOR {
		res += " + "
	}

	// height not shown on approach
	show_height := p.want_height != p.height && p.state != StateAproach
	show_dir := p.want_turn != 0

	if show_height && show_dir {
		res += fmt.Sprintf(" [%s %d]",
			p.Direction.Right(p.want_turn),
			p.want_height)
	} else if show_height {
		res += fmt.Sprintf(" [%d]", p.want_height)
	} else if show_dir {
		res += fmt.Sprintf(" [%s]",
			p.Direction.Right(p.want_turn))
	}
	return res
}

func (p Plane) StateMessage() string {
	res := p.State()
	switch {
	case p.state == StateWaiting:
		res += " -- Awaiting Takeoff --"
	case p.state == StateRolling:
		res += " -- Rolling! --"
	case p.is_holding:
		res += " -- Holding --"
	case p.state == StateAproach:
		res += " -- Final Approach --"
	case p.clear_to_aproach != 0:
		res += " -- Cleared --"
	case p.state == StateLanded:
		res += " -- Landed --"
	case p.state == StateDeparted:
		res += " -- Departed Area --"
	}

	return res
}

func MakePlanes(rules *GameRules, board *Board, diff *Difficulty, seed int64) []*Plane {
	planes := make([]*Plane, 0, diff.num_planes)

	r := rand.New(rand.NewSource(seed))
	plane_types := PlaneTypes(rules)

	for n := 0; n < diff.num_planes; n++ {
		var plane *Plane
		tries := 0

	retry_plane: // try until valid plan found
		for {
			if tries > 100 {
				panic(fmt.Sprintf("cannot find valid plane (found already %d planes)", len(planes)))
			}
			tries++

			typ := ChoosePlaneType(r, plane_types)
			route := ChooseRoute(r, board.routes)

			// entries are present. checked in board.go
			entry := board.entrypoints[route.entry]
			exit := board.entrypoints[route.exit]

			if !typ.entry_exit_routes && !entry.is_airport && !exit.is_airport {
				continue retry_plane
			}

			if !typ.airport_loop && entry == exit && entry.is_airport {
				continue retry_plane
			}

			if !typ.airport_entry && entry.is_airport {
				continue retry_plane
			}

			if !typ.airport_exit && exit.is_airport {
				continue retry_plane
			}

			start := Ticks(RandRange(r, int(rules.last_plane_start), int(diff.duration)))

			height := RandRange(r, typ.entry_min_height, typ.entry_max_height)
			if entry.is_airport {
				height = 0
			}

			plane = &Plane{
				typ: typ,

				entry: entry,
				exit:  exit,

				Position:  entry.Position,
				Direction: route.Direction,

				start:     start,
				fuel_left: typ.initial_fuel,

				height:         height,
				want_height:    height,
				initial_height: height,

				is_holding:   false,
				is_hoovering: typ.can_hoover && entry.is_airport,

				hold_at_navaid: exit.is_airport,
			}

			// no two planes from the same origin share the same altitude<
			for _, other_plane := range planes {
				if other_plane.entry == plane.entry &&
					!other_plane.entry.is_airport &&
					other_plane.initial_height == plane.initial_height &&
					Ticks(Abs(int(other_plane.start-plane.start))) < REUSE_ENTRYPOINT_TIME {
					// retry another plane
					continue retry_plane
				}
			}

			break
		}

		planes = append(planes, plane)
	}

	sort.Sort(ByTime(planes))

	// assign callsigns last because there might not
	// be enough letters for all planes and the first ones
	// should get callsigns first
	callsigns := r.Perm(Min(diff.num_planes, 26))
	for n, callsign := range callsigns {
		planes[n].callsign = rune(callsign + 'A')
	}

	return planes
}

type ByTime []*Plane

func (a ByTime) Len() int           { return len(a) }
func (a ByTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTime) Less(i, j int) bool { return a[i].start > a[j].start }
