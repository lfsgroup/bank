package bank

import "fmt"

type BSB [6]byte

func NewBSB(bsb string) BSB {
	return BSB{}
}

func (b BSB) String() string {
	return fmt.Sprintf("%d%d%d-%d%d%d", b[0], b[1], b[2], b[3], b[4], b[5])
}
