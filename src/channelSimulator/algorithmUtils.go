package channelSimulator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	// "sync"
)

type AlgCfg struct {
	AlgType            string
	VarNodes          uint8
	StateNodes        []uint8
	FuncNodes         uint8
	FuncSpecs         []string
	MessagePassingType string
	Iterations         int
	Decodings          int
	Compute            []string
	Rate               float64
	Eb                 float64
	No                 float64
	Graph              *FactorGraph
}

func cleanInputStrings(in []string) (out []string) {
	out = make([]string, len(in))

	for i := range in {

		out[i] = strings.Fields(in[i])[0]

	}

	return
}

func write_to(filename string, buffer []byte) (err error) {

	f, err := os.Create(filename)
	defer f.Close()
	if err != nil {
		fmt.Printf("Open %s\n", err.Error())
		return
	}
	_, err = f.Write(buffer)

	if err != nil {
		fmt.Printf("Write %s\n", err.Error())
		return
	}
	return nil
}

func toJSON(iface interface{}, im_filename string) (err error) {

	b := new(bytes.Buffer)

	enc := json.NewEncoder(b)
	err = enc.Encode(&iface)

	if err != nil {
		fmt.Printf("encode %s\n", err.Error())
		return
	}

	err = write_to(im_filename+".json", b.Bytes())

	if err != nil {
		fmt.Printf("write %s\n", err.Error())
		return
	}

	return nil
}

func determineDecodingAlg(lines []string) (decoding string, err error) {

	m := 0
	fmt.Sscan(lines[2], &m)
	decoding = strings.Split(lines[3+m], "")[0]

	return
}


func (ac *AlgCfg) Printer(in []chan VariableOut) {

	word := make([]int, ac.VarNodes)
	words := make([][]int,0)
	all := make(chan VariableOut)

	// inWg := new(sync.WaitGroup)
	// once := new(sync.Once)
	for i, ch := range in {
		// inWg.Add(1)
		go func(i int, ch chan VariableOut) {

			for bit := range ch {
				all <- bit
			}
			// inWg.Done()

			// inWg.Wait()
			// once.Do(func() { close(all) })

	}(i, ch)

	}
	i := 0
	for d := range all {
		i++
		word[d.Id] = d.Var
		
		if i % 5 == 0 {
			words = append(words,word)
			print("[ ")
			for i := range word {
				print(word[i], " ")
				word[i] = 0
			}
			print("]\n")
		}





	}
	return
}


func (ac *AlgCfg) FromString(in string) (err error) {

	lines := strings.Split(in, "\n")
	lines = cleanInputStrings(lines)

	ac.AlgType, err = determineDecodingAlg(lines)
	fmt.Sscan(lines[0], &ac.VarNodes)

	if strings.Split(lines[1], ":")[0] == "0" {

	} else {
		tmp := strings.Split(lines[1], ":")[1]
		stateNodes := strings.Split(tmp, ",")
		for i := range stateNodes {
			var j uint8
			fmt.Sscan(stateNodes[i], &j)
			ac.StateNodes = append(ac.StateNodes, j)

		}
	}

	fmt.Sscan(lines[2], &ac.FuncNodes)
	ac.FuncSpecs = lines[3 : 3+ac.FuncNodes]

	switch ac.AlgType {
	case "A", "C":
		if ac.AlgType == "C" {
			fmt.Sscan(lines[4+ac.FuncNodes], &ac.MessagePassingType)
		} else {
			fmt.Sscan(lines[4+ac.FuncNodes], &ac.Iterations)
		}
		fmt.Sscan(lines[5+ac.FuncNodes], &ac.Decodings)
		fmt.Sscan(lines[6+ac.FuncNodes], &ac.Rate)
		fmt.Sscan(lines[7+ac.FuncNodes], &ac.Eb)
		fmt.Sscan(lines[8+ac.FuncNodes], &ac.No)
	case "B":
		ac.Compute = strings.Split(lines[4+ac.FuncNodes], ",")
	default:
		err = fmt.Errorf("No such Algorithm")
		return
	}

	return
}

func (ac *AlgCfg) String() string {

	s := "Algorithm type: " + ac.AlgType + "\n"
	s += "Variable Nodes: "
	s += fmt.Sprint(ac.VarNodes) + "\n"
	s += "State Nodes: "
	s += fmt.Sprint(len(ac.StateNodes)) + " "
	s += fmt.Sprint(ac.StateNodes) + "\n"
	s += "Function Nodes: "
	s += fmt.Sprint(ac.FuncNodes) + "\n"
	s += "Functions: \n\t"
	sa := make([]string, ac.FuncNodes)
	offset := int(ac.VarNodes) + len(ac.StateNodes)
	for i := 0; i < len(sa); i++ {

		sa[i] = "f" + fmt.Sprint(i) + "("
		ed := ac.Graph.Vertices[i+offset].OutEdges
		for j := 0; j < len(ed); j++ {

			sa[i] += fmt.Sprint(ed[j].B.Id)
			if j+1 != len(ed) {
				sa[i] += ","
			}
		}
		sa[i] += ") => "
		sa[i] += fmt.Sprint(ac.Graph.Vertices[i+offset].Output)
	}
	s += strings.Join(sa, "\n\t") + "\n"
	switch ac.AlgType {
	case "A", "C":
		if ac.AlgType == "C" {
			s += "Message-Passing Variant: "
			s += fmt.Sprint(ac.MessagePassingType) + "\n"
		} else {
			s += "Flooding Iterations: "
			s += fmt.Sprint(ac.Iterations) + "\n"
		}
		s += "Decodings: "
		s += fmt.Sprint(ac.Decodings) + "\n"
		s += "Rate: "
		s += fmt.Sprint(ac.Rate) + "\n"
		s += "Eb: "
		s += fmt.Sprint(ac.Eb) + "\n"
		s += "No: "
		s += fmt.Sprint(ac.No) + "\n"
	case "B":
		s += "Compute :"
		s += strings.Join(ac.Compute, ",") + "\n"
	default:

		return "No such Algorithm\n"
	}
	return s
}
