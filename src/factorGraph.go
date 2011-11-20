package channelSimulator

import (
	"errors"
	"sync"
)

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

func (v *Vertex) coms(in []chan T, out []chan T) {
	type msg struct {
		idx int
		// num int
		data T
	}
	on := make([]bool, len(in))
	for i := range on {
		on[i] = true
	}

	all := make(chan msg, 10)

	// if v.Mode == 2 {
	// 	t := T{true, make([]float64, 1)}
	// 	all <- msg{-1, t}
	// }

	inWg := new(sync.WaitGroup)

	once := new(sync.Once)
	// wg.Add(1)
	for i, ch := range in {
		inWg.Add(1)
		go func(i int, ch chan T, id int) {

			for v := range ch {
				println(i, ch, v.H)
				// on[i] = v.H
				if v.H {
					v.P[0] += 1
					all <- msg{i, v}
				} else {
					all <- msg{i, v}
					break
				}

			}
			println(id, "stop listening to", i, ch)
			inWg.Done()
			// on[i] = false
			// if getTrues(on) == 0 {
			inWg.Wait()
			once.Do(func() { close(all) })
			// }

		}(i, ch, v.Id)

	}
	message_number := 0
	outWg := new(sync.WaitGroup)
	for d := range all {

		println("got", v.Id, d.idx)
		if !d.data.First {
			on[d.idx] = false
			message_number++
		}
		tmpCh := onChannels(d.data.First || len(out) == message_number, on, out)

		open := true
		// wg2 := new(sync.WaitGroup)
		for i, ch := range tmpCh {
			println("alli", v.Id, i, d.idx, d.data.String(), len(tmpCh), trues(on))

			if len(out) == message_number {
				outWg.Add(1)
				go func(ch chan T) {
				ch <- T{H: false, P: d.data.P}
				outWg.Done()
				}(ch)
				on[i] = false
				open = false
				continue
			}
			println(ch)
			// wg2.Add(1)
			if d.data.H {
				outWg.Add(1)
				go func(ch chan T) {
					// wg2.Done()
					ch <- T{H: true, P: d.data.P}
					outWg.Done()
				}(ch)
			}

		}
		if !open {
			outWg.Wait()
			closeAll(out)
			break
		}
		// wg.Wait()
		if d.idx >= 0 && d.data.H {
			on[d.idx] = true
		}

	}
	// wg.Done()
}
func (v *Vertex) Run(message chan Monitor, algType string) {

	in := edgeToChannels(v.InEdges)
	out := edgeToChannels(v.OutEdges)
	println(v.Id, "runing?")
	v.coms(in, out)
	println(v.Id, "done?")

}
