package channelSimulator

import (
	// "errors"
	"sync"
)

type VariableOut struct {
	Id  int
	Var int
}
type Message struct {
	From int
	To   int
	Num  int
	Var  []float64
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
	rwLock   sync.RWMutex
	StdOut   chan VariableOut
}

func newVertex(mode uint8, id int) (v *Vertex, err error) {
	v = &Vertex{Mode: mode, Id: id}
	return
}

func (v *Vertex) compressFor(in []float64, x int) (cv []float64) {
	cv = make([]float64, 2)
	for i, c := range in {
		if v.Ttable[i][x] {
			cv[1] += c
		} else {
			cv[0] += c
		}
	}
	return
}

func (v *Vertex) getInputVar(x, i int) (f float64) {
	//  for all that is not x
	f = 1.0
	for j := range v.Fvars {
		if j != x {
			v.rwLock.RLock()
			if v.Ttable[i][j] {
				f *= v.Fvars[j][1]
			} else {
				f *= v.Fvars[j][0]
			}
			v.rwLock.RUnlock()
		}
	}
	return
}

func (v *Vertex) marginOf(x int) (out []float64) {
	expOut := make([]float64, len(v.Output))
	for i := range v.Output {
		inputVar := v.getInputVar(x, i)
		expOut[i] = v.Output[i] * inputVar
	}
	tmpA := v.compressFor(expOut, x)
	out = normalize(tmpA)
	return
}

func (v *Vertex) coms(in []chan T, out []chan T, flood bool) {
	type msg struct {
		idx  int
		data T
	}
	on := make([]bool, len(in))
	for i := range on {
		on[i] = true
	}

	all := make(chan msg, 100)
	if len(out) == 1 {
		switch v.Mode {
		case 0:
			t := T{First: true, H: true, P: normalize(v.Variable)}
			all <- msg{-1, t}
		case 1:
			t := T{First: true, H: true, States: v.States}
			all <- msg{-1, t}
		case 2:
			t := T{First: true, H: true, P: normalize(v.Output)}
			all <- msg{-1, t}
		default:
			println("no such node mode", v.Mode)
			return
		}
	}

	inWg := new(sync.WaitGroup)

	once := new(sync.Once)
	for i, ch := range in {
		inWg.Add(1)
		go func(i int, ch chan T, id int) {
			for v := range ch {
				if v.H {
					all <- msg{i, v}
				} else {
					all <- msg{i, v}
					break
				}

			}
			inWg.Done()
			inWg.Wait()
			once.Do(func() { close(all) })
		}(i, ch, v.Id)
	}
	message_number := 0
	outWg := new(sync.WaitGroup)
	for d := range all {
		if !d.data.First {
			on[d.idx] = false
			message_number++

			switch v.Mode {
			case 0:
				v.rwLock.Lock()
				v.Variable = d.data.P
				v.rwLock.Unlock()

			case 1:
			case 2:
				v.rwLock.Lock()
				if d.data.States != 0 {
					v.Fvars[d.idx] = make([]float64, 1<<d.data.States)
				} else {
					v.Fvars[d.idx] = d.data.P
				}
				v.rwLock.Unlock()

			default:
				return
			}

		}
		var tmpCh []chan T
		if flood {
			tmpCh = onChannels(d.data.First || len(out) == message_number, on, out)
		} else {
			if !((len(out) - 1) <= message_number) {
				continue
			}
			tmpCh = out
		}

		open := true

		for i, ch := range tmpCh {

			if len(out) == message_number {
				outWg.Add(1)
				go func(ch chan T, i int) {
					switch v.Mode {
					case 0:
						ch <- T{H: false, P: v.Variable}
					case 1:
					case 2:
						msg := v.marginOf(i)
						//HACK
						if msg[0] != 1.5 {
							ch <- T{H: false, P: msg}
						}
					}
					outWg.Done()
				}(ch, i)
				open = false
				continue
			}
			if d.data.H {
				if on[i] {
					outWg.Add(1)
					go func(ch chan T, i int) {
						v.rwLock.RLock()
						switch v.Mode {
						case 0:
							ch <- T{H: true, P: v.Variable}
						case 1:
						case 2:
							msg := v.marginOf(i)
							//HACK
							if msg[0] != 1.5 {
								ch <- T{H: true, P: msg}
							}
						}
						v.rwLock.RUnlock()
						outWg.Done()
					}(ch, i)
				}
			}

		}
		if !open {
			outWg.Wait()
			break
		}
		if d.idx >= 0 && flood && d.data.H {
			on[d.idx] = true
		}

	}
}

func (v *Vertex) Init() {
	switch v.Mode {
	case 0:
		v.Variable = []float64{0.5, 0.5}
	case 1:
		// TODO: make states
	case 2:
		v.Ttable = permuteTrues(uint(len(v.OutEdges)))
		v.Fvars = make([][]float64, len(v.OutEdges))
		for i := 0; i < len(v.Fvars); i++ {
			v.Fvars[i] = make([]float64, 2)
		}
	default:
		return
	}

}

func (v *Vertex) Run(T string, decodings int, iterations int, awgn func() (v []float64)) {

	in := edgeToChannels(v.InEdges)
	out := edgeToChannels(v.OutEdges)
	switch T {
	case "A":
		for i := 0; i < decodings; i++ {
			if v.G {
				v.Output = awgn()
			}
			for j := 0; j < iterations; j++ {
				v.coms(in, out, true)
			}
			if v.Mode == 0 {
				v.StdOut <- VariableOut{v.Id, roundFloat64(v.Variable[0])}
			}
		}
		if v.Mode == 0 {
			close(v.StdOut)
		}

	case "B":
		v.coms(in, out, false)
		if v.Mode == 0 {
			println("(", v.Id, ") = {", v.Variable[0], v.Variable[1], "}")
		}
	case "C":
	}
}
