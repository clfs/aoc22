package day10

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
)

type Op struct {
	Name string
	Arg  int
}

func (op *Op) UnmarshalText(text []byte) error {
	fields := strings.Fields(string(text))
	if len(fields) > 2 {
		return fmt.Errorf("invalid op: %q", text)
	}

	switch fields[0] {
	case "noop":
		op.Name = "noop"
		op.Arg = 0
	case "addx":
		var err error
		op.Name = "addx"
		op.Arg, err = strconv.Atoi(fields[1])
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid op: %q", text)
	}

	return nil
}

type Program []Op

func ParseProgram(r io.Reader) (Program, error) {
	var p Program
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		var op Op
		err := op.UnmarshalText(scanner.Bytes())
		if err != nil {
			return nil, err
		}
		p = append(p, op)
	}
	return p, scanner.Err()
}

type CPU struct {
	x  int
	p  Program
	pc int

	executing *Op
	countdown int

	cycle  int
	sprite int
	crt    [][]bool
}

func (c *CPU) Load(p Program) {
	c.x = 1
	c.pc = 0
	c.p = p

	c.executing = nil
	c.countdown = 0

	c.cycle = 1
	c.sprite = 1
	c.crt = make([][]bool, CRTHeight)
	for i := range c.crt {
		c.crt[i] = make([]bool, CRTWidth)
	}
}

// Tick completes one cycle. It returns the value of x during the cycle.
func (c *CPU) Tick() int {
	c.start()
	result := c.during()
	c.after()
	c.cycle += 1
	return result
}

func (c *CPU) start() {
	log.Printf("sprite position: %d", c.sprite)

	// Instructions are only started if we're not waiting for any to complete.
	if c.executing != nil {
		return
	}

	// Start the next instruction.
	c.executing = &c.p[c.pc]

	log.Printf("start cycle %d: begin executing %v", c.cycle, c.executing)

	switch c.executing.Name {
	case "noop":
		c.countdown = 1
	default: // "addx"
		c.countdown = 2
	}
}

func (c *CPU) during() int {
	// Find the pixel we need to draw.
	row := (c.cycle - 1) / CRTWidth
	col := (c.cycle - 1) % CRTWidth

	spriteCol := c.sprite % CRTWidth

	log.Printf("during cycle %d: crt drawing pixel (%d, %d)", c.cycle, row, col)

	if col == spriteCol {
		c.crt[row][col] = true
	} else if col == spriteCol-1 {
		c.crt[row][col] = true
	} else if col == spriteCol+1 {
		c.crt[row][col] = true
	}

	log.Printf("current crt:\n%s", c.Render())

	return c.x
}

func (c *CPU) Render() string {
	var b strings.Builder
	for _, row := range c.crt {
		for _, pixel := range row {
			if pixel {
				b.WriteRune('#')
			} else {
				b.WriteRune('.')
			}
		}
		b.WriteRune('\n')
	}
	return b.String()
}

func (c *CPU) after() {
	// Decrement the countdown.
	c.countdown--

	// If the countdown completed...
	if c.countdown == 0 {
		// Add the addx arg, if any
		if c.executing.Name == "addx" {
			c.x += c.executing.Arg
		}

		log.Printf("end of cycle %d: finish executing %v (Register X is now %d)", c.cycle, c.executing, c.x)

		// Clear the executing instruction.
		c.executing = nil
		// Increment the program counter.
		c.pc++
	}

	c.sprite = c.x
}

func Part1(r io.Reader) (int, error) {
	program, err := ParseProgram(r)
	if err != nil {
		return 0, err
	}

	var cpu CPU
	cpu.Load(program)

	var sum int
	for i := 1; i <= 220; i++ {
		x := cpu.Tick()
		if i%40 == 20 {
			sum += i * x
		}
	}
	return sum, nil
}

const (
	CRTHeight = 6
	CRTWidth  = 40
)

func Part2(r io.Reader) (string, error) {
	program, err := ParseProgram(r)
	if err != nil {
		return "", err
	}

	var cpu CPU
	cpu.Load(program)

	for i := 1; i <= CRTWidth*CRTHeight; i++ {
		cpu.Tick()
	}
	return cpu.Render(), nil
}
