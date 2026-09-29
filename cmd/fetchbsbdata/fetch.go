package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// AusPayNet static download links, the file names stay the same and the content is updated monthly.
// Reference: https://auspaynet.com.au/BSBLinks
const (
	bsbURL         = "https://auspaynetbsbpublic.blob.core.windows.net/bsb-reports/BSBDirectoryFull.csv"
	institutionURL = "https://auspaynetbsbpublic.blob.core.windows.net/bsb-reports/key%20to%20abbreviations%20and%20bsb%20numbers.csv"
)

// bank data files
const (
	bsbFileName             = "data/bsb.csv"
	bsbMetaFileName         = "data/bsbmeta.txt"
	institutionFileName     = "data/institution.csv"
	institutionMetaFileName = "data/institutionmeta.txt"
)

func main() {

	// BSB file download
	bsbFile, lastModified, err := download(bsbURL)
	if err != nil {
		log.Printf("Download error: %v\n", err)
		os.Exit(1)
	}
	err = os.WriteFile(bsbFileName, bsbFile, 0644)
	if err != nil {
		log.Printf("Saving error: %v\n", err)
		os.Exit(1)
	}
	err = saveMetaFile(bsbMetaFileName, bsbURL, lastModified)
	if err != nil {
		log.Printf("Saving error: %v\n", err)
		os.Exit(1)
	}
	log.Printf("Saving %q as %q", bsbURL, bsbFileName)

	// institution file download
	institutionFile, lastModified, err := download(institutionURL)
	if err != nil {
		log.Printf("Download error: %v\n", err)
		os.Exit(1)
	}
	institutionFile, err = convertInstitutions(institutionFile)
	if err != nil {
		log.Printf("Converting institutions error: %v\n", err)
		os.Exit(1)
	}
	err = os.WriteFile(institutionFileName, institutionFile, 0644)
	if err != nil {
		log.Printf("Saving error: %v\n", err)
		os.Exit(1)
	}
	err = saveMetaFile(institutionMetaFileName, institutionURL, lastModified)
	if err != nil {
		log.Printf("Saving error: %v\n", err)
		os.Exit(1)
	}
	log.Printf("Saving %q as %q", institutionURL, institutionFileName)
}

func download(url string) ([]byte, string, error) {
	client := http.Client{Timeout: time.Minute}
	resp, err := client.Get(url)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("get %s: %s", url, resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	if len(bytes.TrimSpace(body)) == 0 {
		return nil, "", fmt.Errorf("get %s: empty file", url)
	}
	return body, resp.Header.Get("Last-Modified"), nil
}

// convertInstitutions converts the AusPayNet key to abbreviations file, which has
// a header and one row per BSB prefix:
//
//	BSB Owner,Owner Name,BSB Prefix
//	ABA,Auswide (a division of MyState Bank Limited),645
//	ABA,Auswide (a division of MyState Bank Limited),656
//
// into one row per institution, with its BSB prefixes comma separated:
//
//	ABA,Auswide (a division of MyState Bank Limited),"645, 656"
func convertInstitutions(file []byte) ([]byte, error) {
	records, err := csv.NewReader(bytes.NewReader(file)).ReadAll()
	if err != nil {
		return nil, err
	}
	if len(records) < 2 {
		return nil, fmt.Errorf("no institutions found")
	}
	header := strings.Join(records[0], ",")
	if !strings.EqualFold(header, "BSB Owner,Owner Name,BSB Prefix") {
		return nil, fmt.Errorf("unexpected header %q", header)
	}

	type institution struct {
		code, name string
		prefixes   []string
	}
	var institutions []*institution
	lookup := make(map[[2]string]*institution)
	for _, rec := range records[1:] {
		if len(rec) < 3 {
			return nil, fmt.Errorf("unexpected record %q", rec)
		}
		code, name, prefix := strings.TrimSpace(rec[0]), strings.TrimSpace(rec[1]), strings.TrimSpace(rec[2])
		key := [2]string{code, name}
		inst, ok := lookup[key]
		if !ok {
			inst = &institution{code: code, name: name}
			lookup[key] = inst
			institutions = append(institutions, inst)
		}
		inst.prefixes = append(inst.prefixes, prefix)
	}

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	for _, inst := range institutions {
		err := w.Write([]string{inst.code, inst.name, strings.Join(inst.prefixes, ", ")})
		if err != nil {
			return nil, err
		}
	}
	w.Flush()
	return buf.Bytes(), w.Error()
}

func saveMetaFile(name, url, lastModified string) error {
	data := fmt.Sprintf("SOURCE_URL = %q\n", url)
	data += fmt.Sprintf("SOURCE_LAST_MODIFIED = %q\n", lastModified)
	data += fmt.Sprintf("LAST_UPDATED = %q\n", time.Now().Format(time.RFC3339))
	err := os.WriteFile(name, []byte(data), 0644)
	if err != nil {
		return fmt.Errorf("writing to %s error: %v", name, err)
	}
	return nil
}
