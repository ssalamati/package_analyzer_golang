package main

import (
    _ "net/http/pprof"
	"bufio"
	"compress/gzip"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
)

// PackageStat holds the package name and the number of files it contains.
type PackageStat struct {
	Name  string
	Files int
}

// ByFiles implements sort.Interface based on the Files field.
type ByFiles []PackageStat

func (a ByFiles) Len() int           { return len(a) }
func (a ByFiles) Less(i, j int) bool { return a[i].Files > a[j].Files }
func (a ByFiles) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func main() {
	// Profiling the code
	go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()

	mirrorURL := flag.String("mirror", "http://ftp.uk.debian.org/debian/dists/stable/main/", "Set the Debian mirror URL")
	topN := flag.Int("top", 10, "Number of top packages to display")

	flag.Parse()

	if len(flag.Args()) < 1 {
		log.Fatalf("Usage: %s <architecture>", os.Args[0])
	}

	architecture := flag.Arg(0)

	contentsURL := fmt.Sprintf("%sContents-%s.gz", *mirrorURL, architecture)

	// Downloading the Contents file
	resp, err := http.Get(contentsURL)
	if err != nil {
		log.Fatalf("Error downloading the Contents file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Error downloading the Contents file: Server returned non-200 status code: %d", resp.StatusCode)
	}

	gzReader, err := gzip.NewReader(resp.Body)
	if err != nil {
		log.Fatalf("Error creating gzip reader: %v", err)
	}
	defer gzReader.Close()

	scanner := bufio.NewScanner(gzReader)
	packageStats := make(map[string]int)

	// Parsing the Contents file
	for scanner.Scan() {
		line := scanner.Text()
		// Assuming the file follows the format: filepath packagenames
		parts := strings.Split(line, " ")
		if len(parts) > 1 {
			// Splitting the package names by comma
			packages := strings.Split(parts[len(parts)-1], ",")
			for _, packageName := range packages {
				packageStats[strings.TrimSpace(packageName)]++
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading Contents file: %v", err)
	}

	stats := make(ByFiles, 0, len(packageStats))
	for name, files := range packageStats {
		stats = append(stats, PackageStat{Name: name, Files: files})
	}

	sort.Sort(stats)

	fmt.Printf("Top %d packages with the most files in %s architecture:\n", *topN, architecture)
	for i := 0; i < *topN && i < len(stats); i++ {
		fmt.Printf("%d. %s - %d files\n", i+1, stats[i].Name, stats[i].Files)
	}

	// Preventing the program from exiting immediately for profiling
	fmt.Println("Press enter to exit...")
    bufio.NewReader(os.Stdin).ReadBytes('\n')
}
