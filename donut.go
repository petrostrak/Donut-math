package main

import (
	"fmt"
	"math"
	"os"
)

const (
	TERM_WIDTH  = 80
	TERM_HEIGHT = 22
)

// N > 0 ? N : 0
// b[o] = ".,-~:;=!*#$@"[N > 0 ? N : 0];
func greaterThanN(N float64) int {
	if N > 0 {
		return int(math.Round(N))
	}
	return 0
}

// k % TERM_WIDTH ? ((char*) b)[k] : '\n'
func rewrite(k, term_width int, b [TERM_HEIGHT][TERM_WIDTH]rune) [TERM_WIDTH]rune {
	if math.Mod(float64(k), float64(term_width)) == 0 {
		return b[k]
	}

	return b['\n']
}

func populateByte(b [TERM_HEIGHT][TERM_WIDTH]rune) {
	for i := 0; i < TERM_HEIGHT; i++ {
		for j := 0; j < TERM_WIDTH; j++ {
			b[i][j] = ' '
		}
	}
}

func populateFloat(z [TERM_HEIGHT][TERM_WIDTH]float64) {
	for i := 0; i < TERM_HEIGHT; i++ {
		for j := 0; j < TERM_WIDTH; j++ {
			z[i][j] = 0
		}
	}
}

func main() {
	var A float64 = 0
	var B float64 = 0
	var i, j float64
	var k int
	z := [TERM_HEIGHT][TERM_WIDTH]float64{}
	b := [TERM_HEIGHT][TERM_WIDTH]rune{}

	fmt.Printf("\x1b[2J")
	for {
		populateByte(b)
		populateFloat(z)
		for j = 0; 2*math.Pi > j; j += 0.07 {
			for i = 0; 2*math.Pi > i; i += 0.02 {
				c := math.Sin(i)
				d := math.Cos(j)
				e := math.Sin(A)
				f := math.Sin(j)
				g := math.Cos(A)
				h := d + 2
				D := 1 / (c*h*e + f*g + 5)
				l := math.Cos(i)
				m := math.Cos(B)
				n := math.Sin(B)
				t := c*h*g - f*e
				x := int(40 + 30*D*(l*h*m-t*n))
				y := int(12 + 15*D*(l*h*n+t*m))
				// o := int(x + 80*y)
				N := 8 * ((f*e-c*d*g)*m - c*d*e - f*g - l*d*n)
				if TERM_HEIGHT > y && y > 0 && x > 0 && TERM_WIDTH > x && D > z[y][x] {
					z[y][x] = D
					// ".,-~:;=!*#$@"[ N > 0 ? N : 0 ]
					b[y][x] = [12]rune{'.', ',', '-', '~', ':', ';', '=', '!', '*', '#', '$', '@'}[greaterThanN(N)]
				}
			}
		}
		fmt.Printf("\x1b[H")
		for k = 0; k < 1761; k++ {
			fmt.Fprintln(os.Stdout, rewrite(k, TERM_WIDTH, b))
			A += 0.00004
			B += 0.00002
		}
	}
}
