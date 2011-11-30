package channelSimulator

import (
	"errors"
	// "sync"
)

type FactorGraph struct {
	Vertices []Vertex
	Edges    []Edge
}

type Vertex struct {
	Mode     uint8 // 0 = var, 1 = state, 2 = func
	Id       int
	States   uint
	Variable []float64
	Fvars    [][]float64
	Ttable   [][]bool
	InEdges  []Edge
	OutEdges []Edge
	Output   []float64
	G        bool
}

type Edge struct {
	A  *Vertex
	B  *Vertex
	Ch chan T
}

type Monitor struct {
	Id    int
	Tag   int
	Value T
}

func NewFactorGraph() (fg *FactorGraph, err error) {

	fg = new(FactorGraph)
	return
}

func newVertex(mode uint8, id int) (v *Vertex, err error) {

	v = &Vertex{Mode: mode, Id: id}

	return
}

func (fg *FactorGraph) AddVertex(mode int) (err error) {
	err = fg.AddVertexState(mode, 0)
	return
}
func (fg *FactorGraph) AddVertexState(mode, state int) (err error) {
	if 0 > mode || mode > 2 {
		err = errors.New("Mode is incorrect")
		return
	}
	id := len(fg.Vertices)
	v, err := newVertex(uint8(mode), id)
	if err != nil {
		return
	}
	fg.Vertices = append(fg.Vertices, *v)

	return
}

func (fg *FactorGraph) AddUndirectedEdge(A, B *Vertex) (err error) {
	if A.Id == B.Id {
		err = errors.New("Cannot edge to self")
		return
	}

	// println("creating undirected", A.Id, B.Id)
	ch1 := make(chan T)

	e := &Edge{A: A, B: B, Ch: ch1}
	fg.Edges = append(fg.Edges, *e)
	A.OutEdges = append(A.OutEdges, *e)
	B.InEdges = append(B.InEdges, *e)

	ch2 := make(chan T)
	e = &Edge{A: B, B: A, Ch: ch2}
	fg.Edges = append(fg.Edges, *e)
	B.OutEdges = append(B.OutEdges, *e)
	A.InEdges = append(A.InEdges, *e)

	return
}

func edgeToChannels(in []Edge) (out []chan T) {
	for _, ch := range in {
		// println("A:", in[i].A.Id, " <" , ch.Ch ,"B: ", in[i].B.Id)
		out = append(out, ch.Ch)
	}

	return
}

func closeAll(in []chan T) {
	for _, ch := range in {
		close(ch)
	}

}

func onChannels(override bool, on []bool, in []chan T) (out []chan T) {

	if override || len(in) == 1 {
		return in
	}
	if len(on) != len(in) {
		return
	}

	for i := range in {
		if on[i] {
			out = append(out, in[i])
		}
	}
	return
}
func normalize(in []float64) (out []float64) {
	d := in[0] + in[1]
	if d != 0.0 {
		return []float64{in[0] / d, in[1] / d}
	}
	return []float64{1.5, 0}

}

