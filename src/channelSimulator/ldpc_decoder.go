package channelSimulator

import "errors"

type LDPC struct {
	blockSize int
}

func NewLDPCDecoder() (ldpc *LDPC, err error) {

	return
}

func (ldpc *LDPC) Decode(in []float64) (out string, err error) {

	if len(in) > ldpc.blockSize {
		err = errors.New("Length of in block is larger than LDPC block size")
		return out, err
	}
	return
}
