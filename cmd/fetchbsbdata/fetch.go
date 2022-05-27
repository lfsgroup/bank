package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jlaffaye/ftp"
)

// FTP details
// Reference: https://bsb.auspaynet.com.au/public/BSB_DB.NSF/0/72E7EB6B4734232ECA2579650017682D/$File/Downloading%20BSB%20Files%20from%20AusPayNet%20via%20FTP.pdf
const (
	ftpAddr      = "bsb.hostedftp.com:21"
	ftpUsername  = "anonymous"
	ftpPassword  = "anonymous"
	ftpDirectory = "~auspaynetftp/BSB"
)

// bank data files
const (
	bsbFileName             = "data/bsb.csv"
	bsbMetaFileName         = "data/bsbmeta.txt"
	institutionFileName     = "data/institution.csv"
	institutionMetaFileName = "data/institutionmeta.txt"
)

func main() {

	var err error

	ftpConn, err := ftpConn(ftpAddr, ftpUsername, ftpPassword, ftpDirectory)
	if err != nil {
		log.Printf("FTP connection error: %v\n", err)
		os.Exit(1)
	}

	// BSS file download
	bsbFile, reportName, err := ftpLatestBSBFile(ftpConn)
	if err != nil {
		log.Printf("FTP error: %v\n", err)
		os.Exit(1)
	}
	err = saveFile(bsbFileName, bsbFile)
	if err != nil {
		log.Printf("Copying error: %v\n", err)
		os.Exit(1)
	}
	err = saveMetaFile(bsbMetaFileName, reportName)
	if err != nil {
		log.Printf("Copying error: %v\n", err)
		os.Exit(1)
	}
	log.Printf("Saving report %q as %q", reportName, bsbFileName)

	// institution file download
	institutionFile, reportName, err := ftpLatestInstitutionFile(ftpConn)
	if err != nil {
		log.Printf("FTP error: %v\n", err)
		os.Exit(1)
	}
	err = saveFile(institutionFileName, institutionFile)
	if err != nil {
		log.Printf("Copying error: %v\n", err)
		os.Exit(1)
	}
	err = saveMetaFile(institutionMetaFileName, reportName)
	if err != nil {
		log.Printf("Copying error: %v\n", err)
		os.Exit(1)
	}
	log.Printf("Saving report %q as %q", reportName, institutionFileName)
}

func saveFile(name string, file io.ReadCloser) error {
	localFile, err := os.Create(name)
	if err != nil {
		return fmt.Errorf("create %s error: %v\n", name, err)
	}
	defer file.Close()
	defer localFile.Close()
	_, err = io.Copy(localFile, file)
	if err != nil {
		return err
	}
	return nil
}

func saveMetaFile(name, latestFileName string) error {
	metaFile, err := os.Create(name)
	if err != nil {
		return fmt.Errorf("create %s error: %v\n", name, err)
	}
	defer metaFile.Close()
	data := fmt.Sprintf("LATEST_FILE_NAME = %q\n", latestFileName)
	data += fmt.Sprintf("LAST_UPDATED = %q\n", time.Now().Format(time.RFC3339))

	_, err = metaFile.WriteString(data)
	if err != nil {
		return fmt.Errorf("writing to %s error: %v\n", name, err)
	}
	return nil
}

func ftpConn(addr, username, password, dir string) (*ftp.ServerConn, error) {
	conn, err := ftp.Dial(addr)
	if err != nil {
		return nil, err
	}
	err = conn.Login(username, password)
	if err != nil {
		return nil, err
	}
	err = conn.ChangeDir(dir)
	if err != nil {
		return nil, err
	}
	log.Printf("Connected to FTP server: %s", addr)
	return conn, nil
}

func ftpLatestBSBFile(conn *ftp.ServerConn) (*ftp.Response, string, error) {
	dir, err := conn.List("")
	if err != nil {
		return nil, "", err
	}
	log.Println("Listing all BSB files")
	var latestReportNum int
	var latestReportFileName string
	for _, file := range dir {
		fileName := file.Name
		if file.Type == ftp.EntryTypeFile &&
			strings.HasPrefix(file.Name, "BSBDirectory") &&
			strings.HasSuffix(file.Name, ".csv") {

			fmt.Println("     " + file.Name)

			fileName = strings.TrimPrefix(fileName, "BSBDirectory")
			fileName = strings.TrimSuffix(fileName, ".csv")
			f := strings.SplitN(fileName, "-", 2)
			if len(f) < 2 {
				continue
			}
			// prevMonth := f[0]
			reportNum, err := strconv.Atoi(f[1])
			if err != nil {
				log.Printf("convert report number %q err: %v", f[1], err)
				continue
			}
			if reportNum > latestReportNum {
				latestReportNum = reportNum
				latestReportFileName = file.Name
			}

		}
	}
	log.Printf("Latest report #%d\n", latestReportNum)
	log.Printf("Downloading report %q", latestReportFileName)
	if latestReportFileName == "" {
		return nil, "", errors.New("could not find latest report")
	}
	resp, err := conn.Retr(latestReportFileName)
	if err != nil {
		return nil, latestReportFileName, err
	}

	return resp, latestReportFileName, nil
}

func ftpLatestInstitutionFile(conn *ftp.ServerConn) (*ftp.Response, string, error) {
	dir, err := conn.List("")
	if err != nil {
		return nil, "", err
	}
	log.Println("Listing all Institution files")
	var latestReportDate time.Time
	var latestReportString string
	var latestReportFileName string
	for _, file := range dir {
		fileName := file.Name
		if file.Type == ftp.EntryTypeFile &&
			strings.HasPrefix(file.Name, "KEY") &&
			strings.HasSuffix(file.Name, ".csv") {

			fmt.Println("     " + file.Name)

			fileName = strings.TrimPrefix(fileName, "KEY TO ABBREVIATIONS AND BSB NUMBERS (")
			fileName = strings.TrimSuffix(fileName, ").csv")

			date, err := time.Parse("January 2006", fileName)
			if err != nil {
				log.Println("time parse err:", err)
				continue
			}
			if date.After(latestReportDate) {
				latestReportDate = date
				latestReportString = fileName
				latestReportFileName = file.Name
			}

		}
	}
	log.Printf("Latest report %s\n", latestReportString)
	log.Printf("Downloading report %q", latestReportFileName)
	if latestReportFileName == "" {
		return nil, "", errors.New("could not find latest report")
	}
	resp, err := conn.Retr(latestReportFileName)
	if err != nil {
		return nil, latestReportFileName, err
	}

	return resp, latestReportFileName, nil
}
