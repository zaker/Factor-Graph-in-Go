package main

import (
	"./channelSimulator"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	// "sync"
)

var fileName *string = flag.String("f", "", "-f \"filename\"")

func main() {
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
	// monchans := make([]chan channelSimulator.Monitor, len(ac.Graph.Vertices))

	// quit := make(chan bool)
	// quitWg := new(sync.WaitGroup)
	// quitOnce := new(sync.Once)
	for i := range ac.Graph.Vertices {
		// quitWg.Add(1)
		// go func(ac *channelSimulator.AlgCfg){
		ac.Graph.Vertices[i].Init()

		// 	quitWg.Done()
		// 	quitWg.Wait()
		// 	quitOnce.Do(func() {quit<-true})
		// }(ac)
	}
	awgn := channelSimulator.RandomAWGNGenerator(ac.Rate, ac.Eb, ac.No)

	catch := make([]chan channelSimulator.VariableOut, ac.VarNodes)
	for i := range ac.Graph.Vertices {
		// quitWg.Add(1)
		// go func(ac *channelSimulator.AlgCfg){

		if ac.Graph.Vertices[i].Mode == 0 {
			catch = append(catch, ac.Graph.Vertices[i].StdOut)
		}
		go ac.Graph.Vertices[i].Run(ac.AlgType, ac.Decodings, ac.Iterations, awgn)

		// 	quitWg.Done()
		// 	quitWg.Wait()
		// 	quitOnce.Do(func() {quit<-true})
		// }(ac)
	}

	if ac.AlgType == "A" {
		go ac.Printer(catch)
	}
	// for i, v := range ac.Graph.Vertices {
	// 	if v.Mode == 2 {
	// 		ac.Graph.Vertices[i].InEdges[0].Ch <- channelSimulator.T{true, true, []float64{1.0}}
	// 		break

	// 	}
	// }

	// for {
	// 	b := false
	select {
	// case b = <-quit:

	// default:

	}
	// 	if b {
	// 		break
	// 	}
	// }

}

func contents(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var result []byte
	buf := make([]byte, 100)
	for {
		n, err := f.Read(buf[0:])
		result = append(result, buf[0:n]...)

		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}
	}
	return string(result), nil
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
