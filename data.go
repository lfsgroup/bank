package bank

import (
	_ "embed"
	"encoding/csv"
	"io"
	"log"
	"strconv"
	"strings"
)

//go:embed data/bsb.csv
var bsbCSV string

//go:embed data/institution.csv
var institutionCSV string

func loadData() map[BSB]Branch {

	institutionLookup := make(map[string][]Institution)
	bsbLookup := make(map[BSB]Branch)

	csvReader := csv.NewReader(strings.NewReader(institutionCSV))
	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("loadData data/institution.csv error: %v", err)
			continue
		}
		institutions := institutionLookup[rec[0]]
		institutions = append(institutions, Institution{
			Code:       rec[0],
			Name:       rec[1],
			BSBNumbers: rec[2],
		})
		institutionLookup[rec[0]] = institutions
	}

	csvReader = csv.NewReader(strings.NewReader(bsbCSV))
	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("loadData data/bsb.csv error: %v", err)
			continue
		}
		bsb, err := NewBSB(rec[0])
		if err != nil {
			log.Printf("loadData data/bsb.csv BSB error: %v", err)
			continue
		}
		bank := findInstitution(bsb, rec[1], institutionLookup[rec[1]])
		bsbLookup[bsb] = Branch{
			BSB:      bsb,
			Name:     rec[2],
			Bank:     bank,
			Address:  rec[3],
			Suburb:   rec[4],
			State:    rec[5],
			Postcode: rec[6],
		}
	}

	return bsbLookup
}

func findInstitution(bsb BSB, name string, banks []Institution) Institution {
	for _, bank := range banks {
		codes := strings.Split(bank.BSBNumbers, ",")
		for _, code := range codes {
			codeNum, err := strconv.Atoi(strings.TrimSpace(code))
			if err != nil {
				log.Printf("findInstitution data/institution.csv matching BSB error: %v", err)
				continue
			}
			var bsbNum int
			if len(code) == 1 || len(code) == 2 {
				bsbNum = int(bsb) / 10000
			} else if len(code) == 3 {
				bsbNum = int(bsb) / 1000
			}
			if codeNum == bsbNum {
				return bank
			}
		}
	}
	return Institution{Code: name}
}
