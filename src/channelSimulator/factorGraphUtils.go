package channelSimulator

import (
	"fmt"
	"strings"
	"math"
)

type T struct {
	H     bool
	First bool
	States uint
	P     []float64
}


// type T interface {
// 	String()
// }
func (t *T) String() string {
	s := ""
	s += fmt.Sprint(t.H)
	s += fmt.Sprint(t.P)
	return s
}

func trues(n []bool) string {
	s := "[ "
	for i := range n {
		s += fmt.Sprint(n[i]) + " "

	}

	return s + "]"
}
func getTrues(n []bool) int {
	t := 0
	for i := range n {
		if n[i] {
			t++
		}

	}

	return t
}

func roundFloat64(f float64) int{
	var t float64
	if f < 0.5 {
		t = math.Floor(f)
	}else {
		t = math.Ceil(f)
	}

	return int(t)

}
func factorialN(n int) int {
	f := 1
	for i := n; i > 1; i-- {
		f *= i
	}
	return f
}

func check(i, c uint) bool {

	chk := uint(1 << uint(i))
	return (c & chk) > 0
}

func int2boolA(n int, l uint) (b []bool) {
	b = make([]bool, 0)
	// print(n," => [ ")
	for i := uint(0); i < l; i++ {

		// println(i,chk)
		// println("n & chk", n & chk)
		b = append(b, check(i, uint(n)))
		// print(b[i], " ")
	}
	// print("]\n")
	return
}

func stringA(inA []float64) string {
	s := fmt.Sprint("[ ")
	for _, in := range inA {
		s += fmt.Sprint(in, ", ")
	}
	s += fmt.Sprint("]")
	return s
}

func stringB(inB []bool) string {
	s := fmt.Sprint("[ ")
	for _, in := range inB {
		s += fmt.Sprint(in, ", ")
	}
	s += fmt.Sprint("]")
	return s
}

// }
func permuteTrues(l uint) [][]bool {
	// numPermutations := factorialN(l)
	numPermutations := 1 << l
	// println(numPermutations)
	trues := make([][]bool, numPermutations)

	for i := 0; i < numPermutations; i++ {
		// b := make([]bool,l)
		// println(stringB(int2boolA(i,uint(l))))
		trues[i] = int2boolA(i, uint(l))
	}
	// println(stringB(trues[0]))


	return trues

}

func MakeGraph(ac *AlgCfg) (err error) {
	g := new(FactorGraph)
	for i := uint8(0); i < ac.VarNodes; i++ {

		g.AddVertex(0)
		g.Vertices[i].StdOut = make(chan VariableOut,4)
	}

	for i := 0; i < len(ac.StateNodes); i++ {

		g.AddVertexState(1, int(ac.StateNodes[i]))

	}
	for i := uint8(0); i < ac.FuncNodes; i++ {

		g.AddVertex(2)

	}
	awgn := RandomAWGNGenerator(ac.Rate, ac.Eb, ac.No)
	for i := 0; i < len(ac.FuncSpecs); i++ {
		fs := strings.Split(ac.FuncSpecs[i], "]")
		input := strings.Split(fs[0][1:], ",")
		separator := ""
		if ac.AlgType == "B" {
			separator += ","
		}

		output := strings.Split(fs[1], separator)
		array := make([]float64, len(output))
		G := false
		for j := 0; j < len(output); j++ {
			if output[j] == "G" {
				G = true
				array = awgn()
				break
			}
			fmt.Sscan(output[j], &array[j])
		}
		node_offset := len(ac.StateNodes) + int(ac.VarNodes)
		g.Vertices[i+node_offset].Output = array
		g.Vertices[i+node_offset].G = G

		for j := 0; j < len(input); j++ {
			vn := 0

			fmt.Sscan(input[j], &vn)

			// fmt.Println(g.Vertices[i+node_offset].Id, "-", g.Vertices[vn].Id)
			err = g.AddUndirectedEdge(&g.Vertices[i+node_offset], &g.Vertices[vn])
			if err != nil {
				return
			}

		}
	}
	ac.Graph = g
	return
}
