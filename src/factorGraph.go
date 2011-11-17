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
	Ch chan float64
}
type T float64

type Monitor struct {
	Id    int
	Tag   int
	Value float64
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
func (fg *FactorGraph) AddDirectedEdge(A, B *Vertex) (e *Edge, err error) {
	if A.Id == B.Id {
		err = errors.New("Cannot edge to self")
		return
	}
	ch := make(chan float64, 4)
	e = &Edge{A: A, B: B, Ch: ch}
	return
}

func (fg *FactorGraph) AddInEdge(A, B *Vertex) (err error) {
	e, err := fg.AddDirectedEdge(A, B)
	fg.Edges = append(fg.Edges, *e)

	A.InEdges = append(A.InEdges, *e)
	B.OutEdges = append(B.OutEdges, *e)
	return
}

func (fg *FactorGraph) AddOutEdge(A, B *Vertex) (err error) {
	e, err := fg.AddDirectedEdge(B, A)
	fg.Edges = append(fg.Edges, *e)

	A.OutEdges = append(A.OutEdges, *e)
	B.InEdges = append(B.InEdges, *e)
	return
}

func (fg *FactorGraph) AddUndirectedEdge(A, B *Vertex) (err error) {
	if A.Id == B.Id {
		err = errors.New("Cannot edge to self")
		return
	}

	println("creating undirected", A.Id,B.Id)
	ch1 := make(chan float64)

	e := &Edge{A: A, B: B, Ch: ch1}
	fg.Edges = append(fg.Edges, *e)
	A.OutEdges = append(A.OutEdges, *e)
	B.InEdges = append(B.InEdges, *e)

	ch2 := make(chan float64)
	e = &Edge{A: B, B: A, Ch: ch2}
	fg.Edges = append(fg.Edges, *e)
	B.OutEdges = append(B.OutEdges, *e)
	A.InEdges = append(A.InEdges, *e)

	return
}


// func vertex(in []chan T, out []chan T) 


func edgeToChannels(in []Edge) (out []chan float64) {
	for i,ch := range in {
		println("A:", in[i].A.Id, " <" , ch.Ch ,"B: ", in[i].B.Id)
		out = append(out, in[i].Ch)
	}

	return
}

func (v *Vertex) monitorChannels(out chan Monitor) {
	in := edgeToChannels(v.InEdges)
	tmp := make(chan Monitor)

	for{
		for i, ch := range in {
			v,ok := <- ch 
			if ok {
				go func(i int, ch chan float64) {
				
					println("FUCKIng",ch,v)
					tmp <- Monitor{i,1, v}
					println("FUCKOff")		
				}(i, ch)
			} else {
				println(ch, "is closed")
			}
			println("done",i)
			out <- <- tmp
		}
		
		println("FU")
	}
	
}

func (v *Vertex) multicastChannels(val float64, exception int,) {
	// you have access to d.idx to know which channel sent the data
	// for _,e := range v.OutEdges {
	// 	println("outing ", val , "to",e.B.Id)
	// }
	out := edgeToChannels(v.OutEdges)
    for i, ch := range out {
    	// if i != exception{
			go func(ch chan float64) {
				println(v.Id,"sending ", val , "on",v.OutEdges[i].Ch)
				ch <- val 
				println(v.Id,"outing ", val , "to",v.OutEdges[i].B.Id)
			}(ch)
		// }
    }
}


func (v *Vertex) coms(in []chan float64, out []chan float64) {
   type msg struct {
    idx  int
    data float64
  }

  all := make(chan msg,10)

  for i, ch := range in {
    go func(i int, ch chan float64,id int) {

      for v := range ch {
      	println(id,"got",v, "on",ch)
        all <- msg{i, v}
      }
    }(i, ch,v.Id)
  }
  for d := range all {
    // you have access to d.idx to know which channel sent the data

    for _, ch := range out {
    	println(v.Id, "sending",d.idx,d.data, "on", ch)
    	go func(ch chan float64){ch <- d.data}(ch)
    }
  }
} 
func (v *Vertex) Run(message chan Monitor, algType string) {

	for _,e := range v.OutEdges {
		println(v.Id,"out ",e.Ch)
	}
	
	for _,e := range v.InEdges {
		println(v.Id,"in ",e.Ch)
	}	

	in := edgeToChannels(v.InEdges)
	out := edgeToChannels(v.OutEdges)

	v.coms(in,out)
	println("done?")
	// inMon := make(chan Monitor,len(v.InEdges)) 
	// inMon := make(chan Monitor) 

	// go v.monitorChannels(inMon)
	// println("gogogo")
	// for {
	// 	// go v.monitorChannels(inMon)
	// 	mon := <- inMon
	// 	val := mon.Value 
	// 	println(v.Id," got ",val, " from ",mon.Id)
	// 	go v.multicastChannels(val,mon.Id)
	// }

}
