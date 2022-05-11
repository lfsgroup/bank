package bank

import (
	"reflect"
	"testing"
)

func TestNewBSB(t *testing.T) {
	tests := []struct {
		bsb     string
		want    BSB
		wantErr bool
	}{
		{bsb: "123-456", want: BSB(123456), wantErr: false},
		{bsb: "012-345", want: BSB(12345), wantErr: false},
		{bsb: "123456", want: BSB(123456), wantErr: false},
		{bsb: "12345", want: BSB(12345), wantErr: false},
		{bsb: "012345", want: BSB(12345), wantErr: false},
		{bsb: "999999", want: BSB(999999), wantErr: false},
		{bsb: "1000000", want: BSB(0), wantErr: true},
		{bsb: "not a number", want: BSB(0), wantErr: true},
		{bsb: "1234", want: BSB(0), wantErr: true},
		{bsb: "0", want: BSB(0), wantErr: true},
		{bsb: "", want: BSB(0), wantErr: true},
		{bsb: "-1", want: BSB(0), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.bsb, func(t *testing.T) {
			got, err := NewBSB(tt.bsb)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewBSB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NewBSB() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBSB_String(t *testing.T) {
	tests := []struct {
		name string
		bsb  BSB
		want string
	}{
		{bsb: BSB(123456), want: "123-456"},
		{bsb: BSB(12345), want: "012-345"},
		{bsb: BSB(999999), want: "999-999"},
		{bsb: BSB(1), want: "000-001"},
		{bsb: BSB(0), want: "000-000"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.bsb.String(); got != tt.want {
				t.Errorf("BSB.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBSB_digits(t *testing.T) {
	tests := []struct {
		name string
		bsb  BSB
		want [6]byte
	}{
		{name: "123-456", bsb: BSB(123456), want: [6]byte{1, 2, 3, 4, 5, 6}},
		{name: "012-345", bsb: BSB(12345), want: [6]byte{0, 1, 2, 3, 4, 5}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.bsb.digits(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BSB.digits(%v) = %v, want %v", tt.bsb, got, tt.want)
			}
		})
	}
}
