package main

import (
	"channelSimulator"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	//	"strings"
)

var fileName *string = flag.String("f", "", "-f \"filename\"")

func main() {
	//TODO: Read config file
	flag.Parse()

	if *fileName == "" {
		fmt.Println("Usage ", os.Args[0], "-f \"filename of config file\"")
		return
	}
	fmt.Println("Reading config")
	fmt.Println(*fileName)
	ac := new(channelSimulator.AlgCfg)
	cfg, _ := contents(*fileName)

	ac.FromString(cfg)
	channelSimulator.MakeGraph(ac)

	fmt.Printf(ac.String())

	fmt.Println("num v ", len(ac.Graph.Vertices))
	monchans := make([]chan channelSimulator.Monitor, len(ac.Graph.Vertices))

	// q := make(chan int) 
	for i := range ac.Graph.Vertices {
		// fmt.Println("Starting", ac.Graph.Vertices[i].Id)
		// go func (ac *channelSimulator.AlgCfg){
		go ac.Graph.Vertices[i].Run(monchans[i], ac.AlgType)
		// q <- 1
		// }(ac)
	}

	for i, v := range ac.Graph.Vertices {
		if v.Mode == 2 {
			ac.Graph.Vertices[i].InEdges[0].Ch <- channelSimulator.T{true, true, []float64{1.0}}
			break

		}
	}

	// i := 0
	select {
	// case j,ok := <- q:
	// i += j
	// println(i)
	// if i >= len(ac.Graph.Vertices) || ok{
	// return
	// }

	}
	println("FU")
}

func contents(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close() // f.Close will run when we're finished.

	var result []byte
	buf := make([]byte, 100)
	for {
		n, err := f.Read(buf[0:])
		result = append(result, buf[0:n]...) // append is discussed later.
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err // f will be closed if we return here.
		}
	}
	return string(result), nil // f will be closed if we return here.
}

func channelSimul() {

	inChannel := make(chan float64, 10)
	outChannel := make(chan float64, 10)
	noiseChannel := make(chan float64, 10)

	noiser := channelSimulator.NewAWGNoise(0.12, 0.0)
	channelModel := channelSimulator.NewChannel(inChannel, noiseChannel, outChannel)

	go messageMaker(inChannel)
	go noiser.ToChannel(noiseChannel)

	go channelModel.Run()

	nm := normCheck()
	for {
		nm(<-outChannel)
	}
}

func messageMaker(in chan float64) {
	for {
		in <- 0.5
	}

}

func normCheck() func(in float64) {
	array := make([]int, 20)

	return func(in float64) {
		index := int(math.Floor(in * 20))
		array[index]++
		fmt.Println(array)
	}
}
