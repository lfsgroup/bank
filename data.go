package bank

import (
	_ "embed"
	"encoding/csv"
	"io"
	"log"
	"strings"
)

//go:embed data/bsb.csv
var bsbCSV string

//go:embed data/institution.csv
var institutionCSV string

func loadData() map[string]Branch {

	institutionLookup := make(map[string]Institution)
	bsbLookup := make(map[string]Branch)

	csvReader := csv.NewReader(strings.NewReader(institutionCSV))
	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err) // TODO: fix this
		}
		institutionLookup[rec[0]] = Institution{
			Code:       rec[0],
			Name:       rec[1],
			BSBNumbers: rec[2],
		}
	}

	csvReader = csv.NewReader(strings.NewReader(bsbCSV))
	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err) // TODO: fix this
		}
		bsbLookup[rec[0]] = Branch{
			BSB:      rec[0],
			Name:     rec[2],
			Bank:     institutionLookup[rec[1]],
			Address:  rec[3],
			Suburb:   rec[4],
			State:    rec[5],
			Postcode: rec[6],
		}
	}

	return bsbLookup
}
