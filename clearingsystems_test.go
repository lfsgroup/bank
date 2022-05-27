package bank

import (
	"testing"
)

func TestNewClearingSystems(t *testing.T) {
	tests := []struct {
		flags string
		want  ClearingSystems
	}{
		// P, PE, PEH, E, EH or blank if closed
		{"", ClosedClearingSystem},
		{"P", PaperClearing},
		{"PE", PaperClearing | ElectronicClearing},
		{"PEH", PaperClearing | ElectronicClearing | HighValueClearing},
		{"E", ElectronicClearing},
		{"EH", ElectronicClearing | HighValueClearing},
		{"P,E,H", PaperClearing | ElectronicClearing | HighValueClearing},
		{"p,e,h", PaperClearing | ElectronicClearing | HighValueClearing},
	}
	for _, tt := range tests {
		t.Run(tt.flags, func(t *testing.T) {
			if got := NewClearingSystems(tt.flags); got != tt.want {
				t.Errorf("NewClearingSystems() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClearingSystems_String(t *testing.T) {
	tests := []struct {
		cs   ClearingSystems
		want string
	}{
		// P, PE, PEH, E, EH or blank if closed
		{ClosedClearingSystem, ""},
		{ClearingSystems(PaperClearing), "P"},
		{ClearingSystems(PaperClearing | ElectronicClearing), "PE"},
		{ClearingSystems(PaperClearing | ElectronicClearing | HighValueClearing), "PEH"},
		{ClearingSystems(ElectronicClearing), "E"},
		{ClearingSystems(ElectronicClearing | HighValueClearing), "EH"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.cs.String(); got != tt.want {
				t.Errorf("ClearingSystems.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClearingSystems_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		cs      ClearingSystems
		want    string
		wantErr bool
	}{
		{
			name: "P,E,H",
			cs:   ClearingSystems(PaperClearing | ElectronicClearing | HighValueClearing),
			want: "\"PEH\"",
		},
		{
			name: "Closed",
			cs:   ClosedClearingSystem,
			want: "\"\"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cs.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("ClearingSystems.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(got) != tt.want {
				t.Errorf("ClearingSystems.MarshalJSON() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestClearingSystems_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		flags   string
		want    ClearingSystems
		wantErr bool
	}{
		{
			name:  "P,E,H",
			flags: "\"P,E,H\"",
			want:  ClearingSystems(PaperClearing | ElectronicClearing | HighValueClearing),
		},
		{
			name:  "PEH",
			flags: "\"PEH\"",
			want:  ClearingSystems(PaperClearing | ElectronicClearing | HighValueClearing),
		},
		{
			name:    "string error",
			flags:   "\"abc",
			want:    ClosedClearingSystem,
			wantErr: true,
		},
		{
			name:    "Closed",
			flags:   "",
			want:    ClosedClearingSystem,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var got ClearingSystems
			err := got.UnmarshalJSON([]byte(tt.flags))
			if (err != nil) != tt.wantErr {
				t.Errorf("ClearingSystems.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("ClearingSystems.UnmarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}
