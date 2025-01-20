package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Point struct {
	x, y int
}

type Cell struct {
	point Point
	cost  int
}

type PriorityQueue []Cell

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool { return pq[i].cost < pq[j].cost }

func (pq PriorityQueue) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }

func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(Cell))
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

func dijkstra(maze [][]int, start, end Point) ([]Point, error) {
	rows := len(maze)
	cols := len(maze[0])
	directions := []Point{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}

	dist := make([][]int, rows)
	prev := make([][]*Point, rows)
	for i := range dist {
		dist[i] = make([]int, cols)
		prev[i] = make([]*Point, cols)
		for j := range dist[i] {
			dist[i][j] = int(^uint(0) >> 1)
		}
	}

	dist[start.x][start.y] = 0
	pq := &PriorityQueue{}
	heap.Init(pq)
	heap.Push(pq, Cell{start, 0})

	for pq.Len() > 0 {
		cell := heap.Pop(pq).(Cell)
		point := cell.point

		if point == end {
			break
		}

		for _, d := range directions {
			neigh := Point{point.x + d.x, point.y + d.y}
			if neigh.x >= 0 && neigh.x < rows && neigh.y >= 0 && neigh.y < cols && maze[neigh.x][neigh.y] != 0 {
				alt := dist[point.x][point.y] + maze[neigh.x][neigh.y]
				if alt < dist[neigh.x][neigh.y] {
					dist[neigh.x][neigh.y] = alt
					prev[neigh.x][neigh.y] = &Point{point.x, point.y}
					heap.Push(pq, Cell{neigh, alt})
				}
			}
		}
	}

	if dist[end.x][end.y] == int(^uint(0)>>1) {
		return nil, fmt.Errorf("path not found")
	}

	path := []Point{}
	for p := &end; p != nil; p = prev[p.x][p.y] {
		path = append([]Point{*p}, path...)
	}

	return path, nil
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	input := []string{}

	line, _ := reader.ReadString('\n')
	input = append(input, strings.TrimSpace(line))

	dimensions := strings.Split(input[0], " ")
	rows, _ := strconv.Atoi(dimensions[0])
	for i := 0; i < rows; i++ {
		line, _ := reader.ReadString('\n')
		input = append(input, strings.TrimSpace(line))
	}

	for i := 0; i < 2; i++ {
		line, _ := reader.ReadString('\n')
		input = append(input, strings.TrimSpace(line))
	}

	size := strings.Split(input[0], " ")
	rows, _ = strconv.Atoi(size[0])
	cols, _ := strconv.Atoi(size[1])

	maze := make([][]int, rows)
	for i := 0; i < rows; i++ {
		maze[i] = make([]int, cols)
		for j, val := range strings.Split(input[1+i], " ") {
			maze[i][j], _ = strconv.Atoi(val)
		}
	}

	startInput := strings.Split(input[1+rows], " ")
	start := Point{atoi(startInput[0]), atoi(startInput[1])}

	endInput := strings.Split(input[2+rows], " ")
	end := Point{atoi(endInput[0]), atoi(endInput[1])}

	path, err := dijkstra(maze, start, end)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for _, p := range path {
		fmt.Printf("%d %d\n", p.x, p.y)
	}
	fmt.Println(".")
}

func atoi(s string) int {
	val, _ := strconv.Atoi(s)
	return val
}
