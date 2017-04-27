package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
)

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	// turn of garbage collection, since we memory isn't an issue here
	debug.SetGCPercent(-1)

	// defer profile.Start(profile.ProfilePath(".")).Stop()
	process(os.Args[1], os.Args[2])
}

var outputFile *os.File
var writer *bufio.Writer

func cleanup() {
	handleErr(writer.Flush())
	handleErr(outputFile.Close())
}
func setupOutput(outputFilname string) {
	var err error
	outputFile, err = os.Create(outputFilname)
	handleErr(err)
	writer = bufio.NewWriter(outputFile)
}

func process(inputFilename, outputFilname string) {
	file, err := os.Open(inputFilename)
	handleErr(err)

	var nodes []node
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if nodes == nil {
			n := parseIndex(line)
			nodes = make([]node, n)
			for i := range nodes {
				nodes[i].index = index(i + 1)
			}
		} else {
			aIndex, bIndex := parseLine(line)
			aNode := &nodes[aIndex-1]
			bNode := &nodes[bIndex-1]
			aNode.addEdge(bNode)
			bNode.addEdge(aNode)
		}
	}
	handleErr(file.Close())
	setupOutput(outputFilname)
	for i := range nodes {
		n := &nodes[i]
		if n.color == white {
			n.color = red
			n.visit()
		}
	}
	foundColoring(nodes)
}

type color uint8

const (
	white color = iota
	red
	blue
)

type index uint

func parseIndex(s string) index {
	i, err := strconv.ParseUint(s, 0, 0)
	handleErr(err)
	return index(i)
}

func parseLine(s string) (index, index) {
	parts := strings.Split(s, " ")
	return parseIndex(parts[0]), parseIndex(parts[1])
}

type node struct {
	color    color
	index    index
	parent   *node
	adjacent []*node
}

// String is for debugging, if we want to print out a node
func (n *node) String() string {
	var parentStr string
	if n.parent != nil {
		parentStr = fmt.Sprint(n.parent.index)
	}
	adjacentIndexes := make([]string, len(n.adjacent))
	for i, aN := range n.adjacent {
		adjacentIndexes[i] = fmt.Sprint(aN.index)
	}
	return fmt.Sprintf("%v %v %v [%v]",
		n.index,
		[]string{"white", "red", "blue"}[n.color],
		parentStr,
		strings.Join(adjacentIndexes, " "),
	)
}

func (n *node) addEdge(anotherN *node) {
	if n.adjacent == nil {
		n.adjacent = []*node{anotherN}
	} else {
		n.adjacent = append(n.adjacent, anotherN)
	}
}

func otherColor(c color) color {
	if c == blue {
		return red
	}
	return blue
}

func (n *node) visit() {
	adjC := otherColor(n.color)
	for _, adjN := range n.adjacent {
		if adjN.color == white {
			adjN.color = adjC
			adjN.parent = n
			adjN.visit()
		} else if adjN.color != adjC {
			foundCycle(adjN, n)
		}
	}
}

func foundCycle(n1, n2 *node) {
	_, err := writer.WriteString("no\n")
	handleErr(err)
	// write the path, starting at n1, up to the root, (but not it)
	for n1.parent != nil {
		writer.WriteString(fmt.Sprintln(n1.index))
		n1 = n1.parent
	}
	// then save the nodes from n2 to the parent, so we can print them in reverse
	path := []*node{n2}
	for n2.parent != nil {
		path = append(path, n2.parent)
		n2 = n2.parent
	}
	// then print them out in reverse order
	for i := len(path) - 1; i >= 0; i-- {
		writer.WriteString(fmt.Sprintln(path[i].index))
	}
	cleanup()
	os.Exit(0)
}

func foundColoring(nodes []node) {
	defer cleanup()
	colors := []string{"white", "red", "blue"}
	_, err := writer.WriteString("yes\n")
	handleErr(err)
	for _, n := range nodes {
		_, err := writer.WriteString(
			fmt.Sprintf("%v %v\n", n.index, colors[n.color]),
		)
		handleErr(err)
	}
}
