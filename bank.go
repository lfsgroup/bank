// bank performs BSB Number lookups to find bank and branch details
//
// Data Source:
// https://bsb.auspaynet.com.au/
package bank

import (
	"errors"
)

var banks data

var ErrBranchNotFound = errors.New("branch not found")

type Institution struct {
	Code       string `json:"code,omitempty"`
	Name       string `json:"name,omitempty"`
	BSBNumbers string `json:"bsb_numbers,omitempty"`
}

type Branch struct {
	BSB      BSB         `json:"bsb,omitempty"`
	Name     string      `json:"name,omitempty"`
	Bank     Institution `json:"bank,omitempty"`
	BankCode string      `json:"bank_code,omitempty"`
	Address  string      `json:"address,omitempty"`
	Suburb   string      `json:"suburb,omitempty"`
	State    string      `json:"state,omitempty"`
	Postcode string      `json:"postcode,omitempty"`
}

func init() {
	banks = loadData()
}

func LookupBSB(bsb string) (Branch, error) {
	return banks.LookupBSB(bsb)
}

type data map[BSB]Branch

func (d data) LookupBSB(bsb string) (Branch, error) {
	bsbNum, err := NewBSB(bsb)
	if err != nil {
		return Branch{}, err
	}
	branch, ok := d[bsbNum]
	if !ok {
		return Branch{}, ErrBranchNotFound
	}
	return branch, nil
}
