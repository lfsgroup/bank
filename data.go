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
	allInstitutions := make([]Institution, 0, 120)
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
		institution := Institution{
			Code:       rec[0],
			Name:       rec[1],
			BSBNumbers: rec[2],
		}
		institutionLookup[rec[0]] = append(institutionLookup[rec[0]], institution)
		allInstitutions = append(allInstitutions, institution)
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
		bank := matchInstitution(bsb, rec[1], institutionLookup[rec[1]], allInstitutions)
		bsbLookup[bsb] = Branch{
			BSB:           bsb,
			Name:          rec[2],
			Bank:          bank,
			BankCode:      rec[1],
			Address:       rec[3],
			Suburb:        rec[4],
			State:         rec[5],
			Postcode:      rec[6],
			PaymentsFlags: NewClearingSystems(rec[7]),
		}
	}

	return bsbLookup
}

func matchInstitution(bsb BSB, name string, subsetBanks []Institution, banks []Institution) Institution {
	bank := matchBSBInstitution(bsb, subsetBanks)
	if bank != (Institution{}) {
		return bank
	}
	bank = matchBSBInstitution(bsb, banks)
	if bank != (Institution{}) {
		return bank
	}
	return Institution{Code: name}
}

func matchBSBInstitution(bsb BSB, banks []Institution) Institution {
	for _, bank := range banks {
		codes := strings.Split(bank.BSBNumbers, ",")
		for _, code := range codes {
			code = strings.TrimSpace(code)
			codeNum, err := strconv.Atoi(code)
			if err != nil {
				log.Printf("findInstitution data/institution.csv matching bsbs[%q] error: %v", code, err)
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
	return Institution{}
}
