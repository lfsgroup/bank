package bank

import (
	"encoding/json"
)

// The payments flags refer to the clearing systems (frameworks) with P=Paper (APCS), E=Electronic (IAC) and H=High Value (HVCS).
type ClearingSystems byte

// ClosedClearingSystem is closed
const ClosedClearingSystem ClearingSystems = 0

const (
	PaperClearing      = 1 << iota // P=Paper (APCS)
	ElectronicClearing             // E=Electronic (IAC)
	HighValueClearing              // H=High Value (HVCS)
)

// NewClearingSystems parses the list of flags into a ClearingSystems type.
// The ClearingsSystem contains all the avaliable clearing systems.
// If no flags are providedd a Closed clearing system is returned.
func NewClearingSystems(flags string) ClearingSystems {
	var cs ClearingSystems
	for _, f := range flags {
		switch f {
		case 'P', 'p':
			cs |= PaperClearing
		case 'E', 'e':
			cs |= ElectronicClearing
		case 'H', 'h':
			cs |= HighValueClearing
		}
	}
	return cs
}

// Closed returns true is the clearing system is closed
func (cs ClearingSystems) Closed() bool { return cs == ClosedClearingSystem }

// Paper returns true if Paper clearing system (APCS) is avaliable
func (cs ClearingSystems) Paper() bool { return cs&PaperClearing == PaperClearing }

// Electronic returns true if Electronic clearing system (IAC) is avaliable
func (cs ClearingSystems) Electronic() bool { return cs&ElectronicClearing == ElectronicClearing }

// HighValue returns true if HighValue clearing system (HVCS) is avaliable
func (cs ClearingSystems) HighValue() bool { return cs&HighValueClearing == HighValueClearing }

// String returns a list of the avaliable clearing systems.
// P=Paper (APCS), E=Electronic (IAC) and H=High Value (HVCS).
// An empty string is returned if closed.
func (cs ClearingSystems) String() string {
	if cs.Closed() {
		return ""
	}
	var str string
	if cs.Paper() {
		str += "P"
	}
	if cs.Electronic() {
		str += "E"
	}
	if cs.HighValue() {
		str += "H"
	}
	return str
}

// MarshalJSON marshals the ClearingSystems into JSON format.
func (cs ClearingSystems) MarshalJSON() ([]byte, error) {
	return json.Marshal(cs.String())
}

// MarshalJSON unmarshal's JSON into a ClearingSystems type.
func (cs *ClearingSystems) UnmarshalJSON(data []byte) error {
	var j string
	err := json.Unmarshal(data, &j)
	if err != nil {
		return err
	}
	*cs = NewClearingSystems(j)
	return nil
}
