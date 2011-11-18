package channelSimulator

import (
	"fmt"
	"strings"
)

type T struct {
	H bool
	P []float64
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

func MakeGraph(ac *AlgCfg) (err error) {
	g := new(FactorGraph)
	for i := uint8(0); i < ac.Var_nodes; i++ {

		g.AddVertex(0)

	}

	for i := 0; i < len(ac.State_nodes); i++ {

		g.AddVertexState(1, int(ac.State_nodes[i]))

	}
	for i := uint8(0); i < ac.Func_nodes; i++ {

		g.AddVertex(2)

	}
	awgn := randomAWGNGenerator(ac.Rate, ac.Eb, ac.No)
	for i := 0; i < len(ac.Func_specs); i++ {
		fs := strings.Split(ac.Func_specs[i], "]")
		input := strings.Split(fs[0][1:], ",")
		separator := ""
		if ac.AlgType == "B" {
			separator += ","
		}

		output := strings.Split(fs[1], separator)
		array := make([]float64, len(output))

		for j := 0; j < len(output); j++ {
			if output[j] == "G" {
				array = awgn()
				break
			}
			fmt.Sscan(output[j], &array[j])
		}
		node_offset := len(ac.State_nodes) + int(ac.Var_nodes)
		g.Vertices[i+node_offset].Output = array

		for j := 0; j < len(input); j++ {
			vn := 0

			fmt.Sscan(input[j], &vn)

			fmt.Println(g.Vertices[i+node_offset].Id, "-", g.Vertices[vn].Id)
			err = g.AddUndirectedEdge(&g.Vertices[i+node_offset], &g.Vertices[vn])
			if err != nil {
				return
			}

		}
	}
	ac.Graph = g
	return
}
