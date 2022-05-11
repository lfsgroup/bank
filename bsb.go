package bank

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// MaxBSB is the maximum number a BSB can hold, as it's only 6 digits.
const MaxBSB = 999999

var ErrInvalidBSB = errors.New("invalid BSB number")

// BSB is an Australian Bank State Branch number.
// It is a code used for domestic payments within Australia.
type BSB uint32

// NewBSB parses a string of the BSB number and returns a valid BSB type.
// If the string is not a valid BSB number an error is returned.
//
// Valid formats:
//
// 		nnn-nnn		example "123-456" will be 123-456
//
// 		nnnnnn 		will convert to nnn-nnn, example "123456" will be 123-456
//
//		nnnnn 		will convert to 0nn-nnn, example "12345" will be 012-345
//
// 	Where n is an ascii digit
//
func NewBSB(bsb string) (BSB, error) {
	var raw [6]byte
	bsb = strings.TrimSpace(bsb)
	if len(bsb) == 7 && bsb[3] == '-' {
		raw[0] = bsb[0]
		raw[1] = bsb[1]
		raw[2] = bsb[2]
		raw[3] = bsb[4]
		raw[4] = bsb[5]
		raw[5] = bsb[6]
	} else if len(bsb) == 6 {
		raw[0] = bsb[0]
		raw[1] = bsb[1]
		raw[2] = bsb[2]
		raw[3] = bsb[3]
		raw[4] = bsb[4]
		raw[5] = bsb[5]
	} else if len(bsb) == 5 {
		raw[0] = '0'
		raw[1] = bsb[0]
		raw[2] = bsb[1]
		raw[3] = bsb[2]
		raw[4] = bsb[3]
		raw[5] = bsb[4]
	} else {
		return 0, ErrInvalidBSB
	}
	num, err := strconv.Atoi(string(raw[:]))
	if err != nil {
		return 0, ErrInvalidBSB
	}
	return intToBSB(num)
}

// MustBSB return a new BSB type.
// Will panic if invalid.
func MustBSB(bsb string) BSB {
	num, err := NewBSB(bsb)
	if err != nil {
		panic(err)
	}
	return num
}

// intToBSB converts an int to a BSB type.
// An ErrInvalidBSB error is return if the bsb is not between 1 - 99999
func intToBSB(bsb int) (BSB, error) {
	if bsb < 1 || bsb > MaxBSB {
		return 0, ErrInvalidBSB
	}
	return BSB(bsb), nil
}

// String return the BSB number in the format nnn-nnn,
// where n is a digit.
// Example:
// 		"123-456"
func (b BSB) String() string {
	bank := b / 1000
	branch := b % 1000
	return fmt.Sprintf("%03d-%03d", bank, branch)
}

// digits return an array for with a byte for each digit.
func (b BSB) digits() [6]byte {
	var bsb [6]byte
	bsb[0] = byte(int(b) / 100000 % 10)
	bsb[1] = byte(int(b) / 10000 % 10)
	bsb[2] = byte(int(b) / 1000 % 10)
	bsb[3] = byte(int(b) / 100 % 10)
	bsb[4] = byte(int(b) / 10 % 10)
	bsb[5] = byte(int(b) / 1 % 10)
	return bsb
}
