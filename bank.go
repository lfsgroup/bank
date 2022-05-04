// bank performs BSB Number lookups to find bank and branch details
//
// Data Source:
// https://bsb.auspaynet.com.au/
package bank

import "errors"

var banks = map[string]Branch{}

type Institution struct {
	Code       string `json:"code,omitempty"`
	Name       string `json:"name,omitempty"`
	BSBNumbers string `json:"bsb_numbers,omitempty"`
}

type Branch struct {
	BSB      string      `json:"bsb,omitempty"`
	Name     string      `json:"name,omitempty"`
	Bank     Institution `json:"bank,omitempty"`
	Address  string      `json:"address,omitempty"`
	Suburb   string      `json:"suburb,omitempty"`
	State    string      `json:"state,omitempty"`
	Postcode string      `json:"postcode,omitempty"`
}

func init() {
	banks = loadData()
}

func LookupBSB(bsb string) (Branch, error) {
	branch, ok := banks[bsb]
	if !ok {
		return Branch{}, errors.New("branch not found")
	}
	return branch, nil
}
