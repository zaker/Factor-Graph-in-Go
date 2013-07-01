package channelSimulator

type Channel struct {
	in    chan float64
	out   chan float64
	noise chan float64
}

func NewChannel(in, noise, out chan float64) (channel *Channel) {

	return &Channel{in: in, out: out, noise: noise}
}

func (ch *Channel) add(chan1, chan2 chan float64) float64 {

	a := <-chan1
	b := <-chan2

	return a + b
}

func (ch *Channel) Run() (err error) {

	print("channel run\n")
	for err == nil {
		ch.out <- ch.add(ch.in, ch.noise)
	}

	println("Loop chrashed ", err.Error())
	return
}
