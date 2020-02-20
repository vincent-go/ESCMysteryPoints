package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	pointsInSpace := make(map[vectorSpace][]point)
	ESCSpeeds := [3]int{1225, 1544, 1864}
	ESCTorques := [3]int{2060, 1850, 1600}
	// ESCSpeeds, ESCTorques := inputSpeedTorques()
	vectorSpaces := create3VectorSpaces(ESCSpeeds, ESCTorques)
	for _, vs := range vectorSpaces {
		pointsInSpace[vs] = pointsInVectorSpace(vs)
	}
	selectedPoints := selectMysteryPoints(pointsInSpace)
	for i := 0; i < 3; i++ {
		fmt.Printf("point selected: %v", selectedPoints[i])
	}
}

type point struct {
	x, y int
}

// String interface impementation for point print
func (p point) String() string {
	return fmt.Sprintf("Speed=%vrpm, Torque=%vNm\n", p.x, p.y)
}

type vectorSpace struct {
	hl point
	hr point
	ll point
	lr point
}

func create3VectorSpaces(ESCSpeeds, ESCTorques [3]int) []vectorSpace {
	p2 := point{x: ESCSpeeds[0], y: ESCTorques[0]}
	p8 := point{x: ESCSpeeds[1], y: ESCTorques[1]}
	p10 := point{x: ESCSpeeds[2], y: ESCTorques[2]}
	p6 := point{x: ESCSpeeds[0], y: int(0.75 * float64(ESCTorques[0]))}
	p4 := point{x: ESCSpeeds[1], y: int(0.75 * float64(ESCTorques[1]))}
	p12 := point{x: ESCSpeeds[2], y: int(0.75 * float64(ESCTorques[2]))}
	p5 := point{x: ESCSpeeds[0], y: int(0.5 * float64(ESCTorques[0]))}
	p3 := point{x: ESCSpeeds[1], y: int(0.5 * float64(ESCTorques[1]))}
	p13 := point{x: ESCSpeeds[2], y: int(0.5 * float64(ESCTorques[2]))}
	p7 := point{x: ESCSpeeds[0], y: int(0.25 * float64(ESCTorques[0]))}
	p9 := point{x: ESCSpeeds[1], y: int(0.25 * float64(ESCTorques[1]))}
	p11 := point{x: ESCSpeeds[2], y: int(0.25 * float64(ESCTorques[2]))}

	vs1 := vectorSpace{
		hl: p2,
		hr: p8,
		ll: p6,
		lr: p4,
	}
	vs2 := vectorSpace{
		hl: p8,
		hr: p10,
		ll: p4,
		lr: p12,
	}
	vs3 := vectorSpace{
		hl: p6,
		hr: p4,
		ll: p5,
		lr: p3,
	}
	vs4 := vectorSpace{
		hl: p4,
		hr: p12,
		ll: p3,
		lr: p13,
	}
	vs5 := vectorSpace{
		hl: p5,
		hr: p3,
		ll: p7,
		lr: p9,
	}
	vs6 := vectorSpace{
		hl: p3,
		hr: p13,
		ll: p9,
		lr: p11,
	}

	vectorSpaces := []vectorSpace{vs1, vs2, vs3, vs4, vs5, vs6}
	vectorSpaces = pick3(vectorSpaces)
	return vectorSpaces
}

func pick3(vectorSpaces []vectorSpace) []vectorSpace {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(vectorSpaces), func(i, j int) { vectorSpaces[i], vectorSpaces[j] = vectorSpaces[j], vectorSpaces[i] })
	return vectorSpaces[:3]
}

// give 10rpm, 10Nm margin for speed and torque
func pointsInVectorSpace(vs vectorSpace) []point {
	var points []point
	slopeLow := float64((vs.lr.y - vs.ll.y) / (vs.lr.x - vs.ll.x))
	slopeHigh := float64((vs.hr.y - vs.hl.y) / (vs.hr.x - vs.hl.x))
	for x := vs.hl.x + 5; x < vs.hr.x-5; x++ {
		yLow := int(slopeLow*float64((x-vs.ll.x))) + vs.ll.y
		yHigh := int(slopeHigh*float64((x-vs.hl.x))) + vs.hl.y
		for y := yLow + 5; y < yHigh-5; y++ {
			point := point{x: x, y: y}
			points = append(points, point)
		}
	}
	return points
}

func selectMysteryPoints(pointsInSpace map[vectorSpace][]point) []point {
	rand.Seed(time.Now().UnixNano())
	var selectedPoints []point
	for _, points := range pointsInSpace {
		selectedPoints = append(selectedPoints, points[rand.Intn(len(points))])
	}
	return selectedPoints
}

func inputSpeedTorques() ([3]int, [3]int) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Please input the A, B, C SPEED, seperate with ','")
	fmt.Print("-> ")
	speedsRawInput, err := reader.ReadString('\n')
	if err != nil {
		log.Panic(err)
	}
	speeds := stringToIntSlice(speedsRawInput)

	fmt.Println("Please input the A, B, C TORQUE, seperate with ','")
	fmt.Print("-> ")
	torquesRawInput, err := reader.ReadString('\n')
	if err != nil {
		log.Panic(err)
	}
	torques := stringToIntSlice(torquesRawInput)
	return speeds, torques
}

func stringToIntSlice(input string) [3]int {
	var output [3]int
	inputAfterSplit := strings.Split(input, ",")
	for i, s := range inputAfterSplit {
		num, err := strconv.Atoi(strings.TrimSpace(s))
		if err != nil {
			log.Panic(err)
		}
		output[i] = num
	}
	return output
}
