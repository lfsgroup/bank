// bank performs BSB Number lookups to find bank and branch details
//
// Data Source:
// https://bsb.auspaynet.com.au/
package bank

import (
	"testing"
)

func TestLookupBSB(t *testing.T) {
	tests := []struct {
		bsb     string
		want    string
		wantErr bool
	}{
		{bsb: "012-020", want: "ANZ PB Project Vantage", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.bsb, func(t *testing.T) {
			got, err := LookupBSB(tt.bsb)
			if (err != nil) != tt.wantErr {
				t.Errorf("LookupBSB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Name != tt.want {
				t.Errorf("LookupBSB() = %v, want %v", got.Name, tt.want)
			}
		})
	}
}
