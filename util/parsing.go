package util

import (
    "encoding/csv"
    "fmt"
    "go/importer"
    "log"
    "os"
)

func ReadCsvFile(filePath string) [][]string {
    f, err := os.Open(filePath)
    if err != nil {
        log.Fatalln("Unable to read input file " + filePath, err)
    }
    defer f.Close()

    csvReader := csv.NewReader(f)
    records, err := csvReader.ReadAll()
    if err != nil {
        log.Fatalln("Unable to parse file as CSV for " + filePath, err)
    }

    return records
}

func DiscoverTypes(packageName string) []string {
    pkg, err := importer.Default().Import(packageName)
    if err != nil {
        log.Fatalln(err)
    }
    return pkg.Scope().Names()
}
