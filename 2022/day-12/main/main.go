// Hill Climbing Algorithm

package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"math"
	"os"
)

type Cell struct {
	x, y      int
	elevation byte
}

type PriorityQueue []*Cell

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].f() < pq[j].f()
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(*Cell))
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]

	return item

}

func (cell *Cell) h(dest *Cell) int {
	return int(math.Abs(float64(cell.x-dest.x)) + math.Abs(float64(cell.y-dest.y)))
}

//
func (cell *Cell) f() int {
	return int(cell.elevation-'a') + 1
}

func findStartAndEnd(grid [][]byte) (*Cell, *Cell) {
	var start, end *Cell
	for x, row := range grid {
		for y, cell := range row {
			if cell == 'S' {
				start = &Cell{x, y, 'a'}
			} else if cell == 'E' {
				end = &Cell{x, y, 'z'}
			}
		}
	}
	return start, end
}

func getNeighbors(grid [][]byte, cell *Cell) []*Cell {
	fmt.Println("Getting neighbors for cell:", cell)

	neighbors := make([]*Cell, 0)
	directions := []struct{ dx, dy int }{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

	for _, dir := range directions {
		newX, newY := cell.x+dir.dx, cell.y+dir.dy

		if newX >= 0 && newX < len(grid) && newY >= 0 && newY < len(grid[0]) {
			neighborElev := grid[newX][newY]
			if int(neighborElev-cell.elevation) <= 1 {
				neighbors = append(neighbors, &Cell{newX, newY, neighborElev})
			}
		}
	}
	return neighbors
}

func calculateAStar(grid [][]byte, start, end *Cell) int {
	openSet := make(PriorityQueue, 0)
	heap.Init(&openSet)
	heap.Push(&openSet, start)

	gScore := make(map[Cell]int)
	fScore := make(map[Cell]int)
	cameFrom := make(map[Cell]*Cell)

	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[0]); j++ {
			gScore[Cell{i, j, 'a'}] = int(math.MaxInt32)
			fScore[Cell{i, j, 'a'}] = int(math.MaxInt32)
		}
	}

	gScore[*start] = 0
	fScore[*start] = start.h(end)

	for len(openSet) > 0 {
		current := heap.Pop(&openSet).(*Cell)

		if *current == *end {
			return gScore[*current]
		}

		for _, neighbor := range getNeighbors(grid, current) {
			tentativeGScore := gScore[*current] + neighbor.f()

			if tentativeGScore < gScore[*neighbor] || gScore[*neighbor] == 0 {
				cameFrom[*neighbor] = current
				gScore[*neighbor] = tentativeGScore
				fScore[*neighbor] = tentativeGScore + neighbor.h(end)

				if !contains(openSet, neighbor) {
					heap.Push(&openSet, neighbor)
				}
			}
		}
	}

	return -1
}

func contains(cells []*Cell, cell *Cell) bool {
	for _, c := range cells {
		if *c == *cell {
			return true
		}
	}
	return false
}

func printGrid(grid [][]byte) {
	for _, row := range grid {
		fmt.Println(string(row))
	}
}

func main() {
	fmt.Println("Advent of Code - 2022 - Day 12 - Hill Climbing Algorithmn")

	data, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening input file:", err)
		return
	}
	defer data.Close()

	scanner := bufio.NewScanner(data)

	var grid [][]byte

	for scanner.Scan() {
		grid = append(grid, []byte(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	start, end := findStartAndEnd(grid)
	steps := calculateAStar(grid, start, end)

	fmt.Println("Fewest steps: ", steps)
}
