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

	if *fileName == ""{
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
	for i := range ac.Graph.Vertices {
		fmt.Println("Starting", ac.Graph.Vertices[i].Id)
		go ac.Graph.Vertices[i].Run(monchans[i], ac.AlgType)
	}

	// for i, ch := range monchans {
	// 	println("unlock", i)
		
	// }

	// ac.Graph.Vertices[0].InEdges[0].Ch <- 0.1
	// ac.Graph.Vertices[0].OutEdges[0].Ch <- 0.2
	// ac.Graph.Vertices[0].InEdges[0].Ch <- 0.1
	// ac.Graph.Vertices[0].OutEdges[0].Ch <- 0.2

	// ac.Graph.Vertices[0].InEdges[0].Ch <- 0.1
	// ac.Graph.Vertices[1].InEdges[0].Ch <- 0.2
	// ac.Graph.Vertices[0].InEdges[0].Ch <- 0.1
	// ac.Graph.Vertices[1].InEdges[0].Ch <- 0.2
	// edges := ac.Graph.Edges

	// for _,e := range edges {
		
	// 	println(e.A.Id, " to ", e.B.Id)
	// }
	go func(){
		in := 0.1
		for i := 0; i < 100;i++{ 
			ac.Graph.Vertices[0].InEdges[0].Ch <- in

			in += 0.1
	// ac.Graph.Vertices[0].OutEdges[0].Ch <- 0.2
			// ac.Graph.Vertices[0].InEdges[0].Ch <- 0.1
	// ac.Graph.Vertices[1].InEdges[0].Ch <- 0.2
		// ac.Graph.Vertices[0].InEdges[0].Ch <- 0.3
		// ac.Graph.Vertices[0].InEdges[0].Ch <- 0.3
		}
	}()
	for{}
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
