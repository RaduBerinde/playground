package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var input = flag.String("input", "", "Input tabular file")
var table = flag.String("table", "", "database.table")
var rowsPerInsert = flag.Int("rows", 50, "rows per insert statement")

func main() {
	flag.Parse()
	if *input == "" {
		log.Fatalf("Specify input file with -input")
	}
	if *table == "" {
		log.Fatalf("Specify input file with -input")
	}
	fi, err := os.Open(*input)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(fi)
	read := func() []string {
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}
			return nil
		}
		return strings.Split(scanner.Text(), "\t")
	}

	// read the headers
	hdrs := read()
	if hdrs == nil {
		log.Fatal("empty file")
	}
	for i := range hdrs {
		// Trim byte order marks.
		hdrs[i] = string(bytes.Trim([]byte(hdrs[i]), "\xef\xbb\xbf"))
	}
	intCols := map[string]bool{
		"taxonID":             true,
		"acceptedNameUsageID": true,
		"parentNameUsageID":   true,
	}

	flush := func(rows [][]string) {
		if len(rows) == 0 {
			return
		}
		fmt.Printf("INSERT INTO %s VALUES\n", *table)
		for rowIdx, r := range rows {
			fmt.Printf("  (")
			for i, val := range r {
				if i > 0 {
					fmt.Printf(",")
				}
				if val == "" {
					fmt.Printf("NULL")
				} else {
					val = strings.Replace(val, "'", "", -1)

					if intCols[hdrs[i]] {
						fmt.Printf("%s", val)
					} else {
						fmt.Printf("'%s'", val)
					}
				}
			}
			if rowIdx == len(rows)-1 {
				fmt.Printf(");\n")
			} else {
				fmt.Printf("),\n")
			}
		}
	}

	var rows [][]string

	n := 0
	for {
		row := read()
		if row == nil {
			break
		}
		rows = append(rows, row)
		if len(rows) >= *rowsPerInsert {
			flush(rows)
			rows = rows[:0]
		}
		n++
		if n%1000 == 0 {
			fmt.Fprintf(os.Stderr, "%dk rows done\n", n/1000)
		}
	}
	flush(rows)
}
