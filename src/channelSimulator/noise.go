package channelSimulator

import (
	"math"
	"math/rand"
	"time"
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
		out <- rand.NormFloat64()*awgn.std_deviation + awgn.mean
	}
}

func RandomAWGNGenerator(rate, eb, no float64) func() (v []float64) {

	std_dev := math.Sqrt(no / 2)
	awgn := NewAWGNoise(std_dev, 0.0)
	rand.Seed(int64(time.Now().Nanosecond()))
	return func() (v []float64) {
		ec := eb * rate
		ti := math.Sqrt(ec) * -1.0
		vi := rand.NormFloat64()*awgn.std_deviation + awgn.mean
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
