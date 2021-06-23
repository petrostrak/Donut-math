package main

import (
	"fmt"
	"math"
)

const (
	THETA_SPACING = 0.07
	PHI_SPACING   = 0.02
	SCREEN_WIDTH  = 1440
	SCREEN_HEIGHT = 900
	R1            = 1
	R2            = 2
	K2            = 5
	// Calculate K1 based on screen size: the maximum x-distance occurs
	// roughly at the edge of the torus, which is at x=R1+R2, z=0.  we
	// want that to be displaced 3/8ths of the width of the screen, which
	// is 3/4th of the way from the center to the side of the screen.
	// screen_width*3/8 = K1*(R1+R2)/(K2+0)
	// screen_width*K2*3/(8*(R1+R2)) = K1
	K1 = SCREEN_WIDTH * K2 * 3 / (8 * (R1 + R2))
)

func renderFrame(A, B float64) {
	// precompute sines and cosines of a and b
	cosA := math.Cos(A)
	cosB := math.Cos(B)
	sinA := math.Sin(A)
	sinB := math.Sin(B)

	output := [][]byte{make([]byte, SCREEN_WIDTH), make([]byte, SCREEN_HEIGHT)}
	zbuffer := [][]float64{{SCREEN_WIDTH}, {SCREEN_HEIGHT}}

	// theta goes around the cross-sectional circle of a torus
	var theta float64
	for theta = 0; theta <= math.Pi; theta += THETA_SPACING {
		// precompute sines and cosines of theta
		cosTheta := math.Cos(theta)
		sinTheta := math.Sin(theta)

		var phi float64
		for phi = 0; phi < 2*math.Pi; phi += PHI_SPACING {
			// precompute sines and cosines of phi
			cosPhi := math.Cos(phi)
			sinPhi := math.Sin(phi)

			// the x,y coordinate of the circle, before revolving (factored
			// out of the above equation
			circlex := R2 + R1*cosTheta
			circley := R1 * sinTheta

			// final 3D (x,y,z) coordinate after rotations, directly from
			// our math above
			x := circlex*(cosB*cosPhi+sinA*sinB*sinPhi) - circley*cosA*sinB
			y := circlex*(sinB*cosPhi-sinA*cosB*sinPhi) + circley*cosA*cosB
			z := K2 + cosA*circlex*sinPhi + circley*sinA
			// "one over z"
			ooz := 1 / z

			// x and y projection.  note that y is negated here, because y
			// goes up in 3D space but down on 2D displays.
			xp := int(SCREEN_WIDTH/2 + K1*ooz*x)
			yp := int(SCREEN_HEIGHT/2 - K1*ooz*y)

			// calculate luminance.  ugly, but correct.
			L := cosPhi*cosTheta*sinB - cosA*cosTheta*sinPhi - sinA*sinTheta + cosB*(cosA*sinTheta-cosTheta*sinA*sinPhi)

			// L ranges from -sqrt(2) to +sqrt(2).  If it's < 0, the surface
			// is pointing away from us, so we won't bother trying to plot it.
			if L > 0 {
				// test against the z-buffer.  larger 1/z means the pixel is
				// closer to the viewer than what's already plotted.
				if ooz > zbuffer[xp][yp] {
					zbuffer[xp][yp] = ooz
					luminance_index := int(L * 8)
					// luminance_index is now in the range 0..11 (8*sqrt(2) = 11.3)
					// now we lookup the character corresponding to the
					// luminance and plot it in our output:
					output[xp][yp] = ".,-~:;=!*#$@"[luminance_index]
				}
			}
		}
	}

	// now, dump output[] to the screen.
	// bring cursor to "home" location, in just about any currently-used
	// terminal emulation mode
	fmt.Printf("\x1b[H")
	for j := 0; j < SCREEN_HEIGHT; j++ {
		for i := 0; i < SCREEN_WIDTH; i++ {
			fmt.Print(output[i][j])
		}
		fmt.Println()
	}
}

func main() {
	renderFrame(1, 2)
}
