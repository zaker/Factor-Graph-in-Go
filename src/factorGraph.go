package channelSimulator

import "errors"

type FactorGraph struct {
	Vertices []Vertex
	Edges    []Edge
}

type Vertex struct {
	Mode     uint8 // 0 = var, 1 = state, 2 = func
	Id       int
	States   int
	InEdges  []Edge
	OutEdges []Edge
	Output   []float64
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
		err = errors.New("Mode is not correct")
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

	println("creating undirected", A.Id, B.Id)
	ch1 := make(chan T,4)

	e := &Edge{A: A, B: B, Ch: ch1}
	fg.Edges = append(fg.Edges, *e)
	A.OutEdges = append(A.OutEdges, *e)
	B.InEdges = append(B.InEdges, *e)

	ch2 := make(chan T,4)
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


func (v *Vertex) coms(in []chan T, out []chan T) {
	type msg struct {
		idx  int
		data T
	}
	n := make([]bool, len(in))
	for i := range n {
		n[i] = true
	}

	all := make(chan msg, 10)
	if v.Mode == 2 {
		t := T{true, make([]float64, 1)}
		all <- msg{-1, t}
	}

	for i, ch := range in {
		go func(i int, ch chan T, id int) {

			for v := range ch {

				v.P[0] += 0.2
				all <- msg{i, v}

			}
			println(id, "stop listening to", i, ch)

			n[i] = false
			if getTrues(n) == 0 {
				close(all)
			}

		}(i, ch, v.Id)
	
	}

	for d := range all {

		// you have access to d.idx to know which channel sent the data

			println("all", v.Id, d.idx, d.data.String(), len(out), trues(n))

			for i, ch := range out {
				println("alli", v.Id, i, d.idx, d.data.String(), len(out), trues(n))
				if d.idx == -1 {

					go func(ch chan T) {
						ch <- T{true, d.data.P}
					}(ch)
					continue
				}
				if getTrues(n) == 0  && v.Mode == 2{
					close(out[i])
					close(in[i])
					return
				}
				if len(out) == 1 {
					println("one left", v.Id, i, d.idx, d.data.String(), len(out), trues(n))
					n[d.idx] = false
					ch <- T{false, d.data.P}
					close(out[i])
					close(in[i])

					return
				}
				if !d.data.H || !n[i] {
					println("blocked", v.Id, i, d.idx, d.data.String(), len(out), trues(n))


					n[d.idx] = false

					continue
				}
				if i != d.idx {

					println(v.Id, i, d.idx, d.data.String(), len(out), trues(n))
					go func(ch chan T) {
						ch <- T{true, d.data.P}
					}(ch)
				}
			}

	}
}
func (v *Vertex) Run(message chan Monitor, algType string) {

	in := edgeToChannels(v.InEdges)
	out := edgeToChannels(v.OutEdges)

	v.coms(in, out)
	println(v.Id, "done?")

}
