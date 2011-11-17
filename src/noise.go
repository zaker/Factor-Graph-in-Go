package channelSimulator

import (
	"math"
	"math/rand"
)

type AWGNoise struct {
	std_deviation float64
	mean          float64
}

func NewAWGNoise(std_deviation, mean float64) *AWGNoise {
	awgn := new(AWGNoise)
	awgn.std_deviation = std_deviation
	awgn.mean = mean

	return awgn

}
func (awgn *AWGNoise) ToChannel(out chan float64) {
	for {
		//		r := rand.NormFloat64()*awgn.std_deviation + awgn.mean
		//		out <- math.Abs(math.Sin(r))
		out <- rand.NormFloat64()*awgn.std_deviation + awgn.mean
	}
}

func randomAWGNGenerator(rate, eb, no float64) func() (v []float64) {
	//	println(rate, eb, no)
	noiseToRand := make(chan float64, 4)

	std_dev := math.Sqrt(no / 2)
	noiser := NewAWGNoise(std_dev, 0.0)

	go noiser.ToChannel(noiseToRand)
	return func() (v []float64) {
		ec := eb * rate
		//		snr := eb /no
		ti := math.Sqrt(ec) * -1.0
		vi := <-noiseToRand
		ri := ti + vi
		divisor := 1.0
		divisor += math.Pow(math.E, -2.0*((math.Sqrt(ec)*ri)/(math.Pow(std_dev, 2))))

		p := 1 / divisor
		v = append(v, p)
		v = append(v, 1.0-v[0])
		v[0], v[1] = v[1], v[0]
		return
	}

}
