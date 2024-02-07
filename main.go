package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// City represents a city with coordinates
type City struct {
	X float64
	Y float64
}

// Distance calculates the Euclidean distance between two cities
func (c *City) Distance(other *City) float64 {
	dx := c.X - other.X
	dy := c.Y - other.Y
	return math.Sqrt(dx*dx + dy*dy)
}

// Ant represents an ant agent
type Ant struct {
	Tour    []int
	Visited map[int]bool
}

// AntColony represents an ant colony
type AntColony struct {
	NumAnts        int
	Alpha          float64
	Beta           float64
	Rho            float64
	Q              float64
	Cities         []*City
	Pheromones     [][]float64
	DistanceMatrix [][]float64
}

// NewAntColony initializes a new ant colony
func NewAntColony(numAnts int, alpha, beta, rho, q float64, cities []*City) *AntColony {
	colony := &AntColony{
		NumAnts:        numAnts,
		Alpha:          alpha,
		Beta:           beta,
		Rho:            rho,
		Q:              q,
		Cities:         cities,
		Pheromones:     make([][]float64, len(cities)),
		DistanceMatrix: make([][]float64, len(cities)),
	}
	for i := range colony.Pheromones {
		colony.Pheromones[i] = make([]float64, len(cities))
	}
	for i := range colony.DistanceMatrix {
		colony.DistanceMatrix[i] = make([]float64, len(cities))
		for j := range colony.DistanceMatrix[i] {
			colony.DistanceMatrix[i][j] = cities[i].Distance(cities[j])
		}
	}
	return colony
}

// InitializeAnts initializes ants with random starting cities
func (ac *AntColony) InitializeAnts() []*Ant {
	ants := make([]*Ant, ac.NumAnts)
	for i := range ants {
		ants[i] = &Ant{
			Tour:    make([]int, len(ac.Cities)),
			Visited: make(map[int]bool),
		}
		startCity := rand.Intn(len(ac.Cities))
		ants[i].Tour[0] = startCity
		ants[i].Visited[startCity] = true
	}
	return ants
}

// NextCity selects the next city for an ant to visit based on pheromone trails and heuristic information
func (ac *AntColony) NextCity(ant *Ant) int {
	currentCity := ant.Tour[len(ant.Tour)-1]
	pheromones := ac.Pheromones[currentCity]
	heuristic := make([]float64, len(ac.Cities))
	sum := 0.0
	for i, city := range ac.Cities {
		if !ant.Visited[i] {
			heuristic[i] = 1 / ac.DistanceMatrix[currentCity][i]
			sum += math.Pow(pheromones[i], ac.Alpha) * math.Pow(heuristic[i], ac.Beta)
		}
	}
	roulette := rand.Float64() * sum
	cumulativeProbability := 0.0
	for i, city := range ac.Cities {
		if !ant.Visited[i] {
			cumulativeProbability += math.Pow(pheromones[i], ac.Alpha) * math.Pow(heuristic[i], ac.Beta)
			if cumulativeProbability >= roulette {
				return i
			}
		}
	}
	// This should not happen
	return -1
}

// AntsMove performs the movement of all ants
func (ac *AntColony) AntsMove(ants []*Ant) {
	for _, ant := range ants {
		for len(ant.Tour) < len(ac.Cities) {
			nextCity := ac.NextCity(ant)
			ant.Tour = append(ant.Tour, nextCity)
			ant.Visited[nextCity] = true
		}
	}
}

// UpdatePheromones updates the pheromone trails based on the tours of the ants
func (ac *AntColony) UpdatePheromones(ants []*Ant) {
	for i := range ac.Pheromones {
		for j := range ac.Pheromones[i] {
			ac.Pheromones[i][j] *= (1 - ac.Rho)
		}
	}
	for _, ant := range ants {
		tourLength := ac.TourLength(ant.Tour)
		for i := 0; i < len(ant.Tour)-1; i++ {
			fromCity := ant.Tour[i]
			toCity := ant.Tour[i+1]
			ac.Pheromones[fromCity][toCity] += ac.Q / tourLength
			ac.Pheromones[toCity][fromCity] += ac.Q / tourLength
		}
	}
}

// TourLength calculates the total length of a tour
func (ac *AntColony) TourLength(tour []int) float64 {
	length := 0.0
	for i := 0; i < len(tour)-1; i++ {
		fromCity := tour[i]
		toCity := tour[i+1]
		length += ac.DistanceMatrix[fromCity][toCity]
	}
	return length
}

func main() {
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Create cities
	cities := []*City{
		{X: 0, Y: 0},
		{X: 1, Y: 1},
		{X: 2, Y: 2},
		{X: 3, Y: 3},
		{X: 4, Y: 4},
	}

	// Set ACO parameters
	numAnts := 10
	alpha := 1.0
	beta := 2.0
	rho := 0.5
	q := 100.0

	// Create ant colony
	colony := NewAntColony(numAnts, alpha, beta, rho, q, cities)

	// Run ACO algorithm
	iterations := 100
	for i := 0; i < iterations; i++ {
		ants := colony.InitializeAnts()
		colony.AntsMove(ants)
		colony.UpdatePheromones(ants)
	}

	// Find best tour
	bestTour := make([]int, len(cities))
	bestTourLength := math.Inf(1)
	for _, ant := range colony.InitializeAnts() {
		tourLength := colony.TourLength(ant.Tour)
		if tourLength < bestTourLength {
			bestTourLength = tourLength
			copy(bestTour, ant.Tour)
		}
	}

	// Print results
	fmt.Println("Best tour:", bestTour)
	fmt.Println("Best tour length:", bestTourLength)
}
