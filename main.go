package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/jung-kurt/gofpdf"
)

func unite(directory string, outfile string, pages []int) string {
	a := []string{}
	for _, p := range pages {
		a = append(a, fmt.Sprintf("%v/page%d.pdf", directory, p))
	}
	a = append(a, fmt.Sprintf("%v/%v", directory, outfile))
	cmd := exec.Command("pdfunite", a...)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	return outfile
}

func frontPages(src []int) []int {
	arr := make([]int, 0, len(src)/2)
	for index, v := range src {
		if index%4 == 0 || index%4 == 1 {
			arr = append(arr, v)
		}
	}
	return arr
}
func backPages(src []int) []int {
	arr := make([]int, 0, len(src)/2)
	for index, v := range src {
		if index%4 == 2 || index%4 == 3 {
			arr = append(arr, v)
		}
	}
	return arr
}

func getPages(start int, end int, total int) []int {
	arr := make([]int, total)
	// n, 1, 2, n-1, n-2, 3, 4, n-3, n-4, 5, 6, n-5, n-6, 7, 8, n-7, n-8, 9, 10, n-9, n-10, 11, 12, n-11…
	offs := 1
	arr[0] = total
	rev := total - 1
	for i := 1; i <= total/2; i += 2 {
		// fmt.Printf("%d out of %d\n", offs, total)
		arr[offs] = i
		offs++
		arr[offs] = i + 1
		offs++
		arr[offs] = rev
		offs++
		rev--
		if offs >= total {
			break
		}
		arr[offs] = rev
		offs++
		rev--
	}

	return arr[start:end]
}

func clean(directory string) {
	dirRead, _ := os.Open(directory)
	dirFiles, _ := dirRead.Readdir(0)

	for index := range dirFiles {
		fileHere := dirFiles[index]
		nameHere := fileHere.Name()
		if strings.Contains(nameHere, ".pdf") && strings.Contains(nameHere, "page") {
			fullPath := directory + "/" + nameHere
			os.Remove(fullPath)
		}
	}
}

func alignedNumOfPages(n int, groupSize int) int {
	if n%groupSize == 0 {
		return n
	}
	return (n/groupSize)*groupSize + groupSize
}

func actualNumOfPages(directory string) int {
	totals := 0
	dirRead, _ := os.Open(directory)
	dirFiles, _ := dirRead.Readdir(0)
	for index := range dirFiles {
		fileHere := dirFiles[index]
		nameHere := fileHere.Name()
		if strings.Contains(nameHere, "page") && strings.Contains(nameHere, ".pdf") {
			totals++
		}
	}
	return totals
}

func createEmptyPage(directory string, pageNum int) {
	pdf := gofpdf.New("P", "mm", "Letter", "")
	pdf.AddPage()
	pdf.Ln(20)
	if err := pdf.OutputFileAndClose(fmt.Sprintf("%v/page%d.pdf", directory, pageNum)); err != nil {
		log.Fatal("pdf creation error: ", err)
	}
}

var nGroup = flag.Int("n", 8, "number of pages in the group. Only 4,8,12,16,20,24,28,32 are allowed")
var sInput = flag.String("i", "input.pdf", "name of the input file")

func main() {
	fmt.Println("PDFBOOKLET")
	flag.Parse()

	if _, err := os.Stat(*sInput); os.IsNotExist(err) {
		log.Fatal("input file not found")
	}

	validGroups := map[int]int{4: 1, 8: 1, 12: 1, 16: 1, 20: 1, 24: 1, 28: 1, 32: 1}
	if _, found := validGroups[*nGroup]; !found {
		log.Fatal("invalid group size")
	}
	pdfseparate, err := exec.LookPath("pdfseparate")
	if err != nil {
		log.Fatal("pdfseparate was not found")
	}
	fmt.Printf("pdfseparate is available at %s\n", pdfseparate)

	pdfunite, err := exec.LookPath("pdfunite")
	if err != nil {
		log.Fatal("pdfunite was not found")
	}
	fmt.Printf("pdfunite is available at %s\n", pdfunite)

	fmt.Printf("input = %s, group by %d\n", *sInput, *nGroup)
	outDir := "."
	clean(outDir)

	cmd := exec.Command("pdfseparate", *sInput, "page%d.pdf")
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	pageNum := actualNumOfPages(outDir)
	totalPageNum := alignedNumOfPages(pageNum, *nGroup)
	fmt.Printf("num pages=%d, aligned=%d\n", pageNum, totalPageNum)
	for n := pageNum + 1; n <= totalPageNum; n++ {
		createEmptyPage(outDir, n)
	}
	for gstart := 1; gstart < totalPageNum; gstart += *nGroup {
		gend := gstart + *nGroup - 1
		pages := getPages(gstart-1, gend, totalPageNum)
		frontFile := fmt.Sprintf("output-%04d-%04d-front.pdf", gstart, gend)
		backFile := fmt.Sprintf("output-%04d-%04d-back.pdf", gstart, gend)
		fmt.Printf("lp -o number-up=2 %v # %v\n", unite(outDir, frontFile, frontPages(pages)), frontPages(pages))
		fmt.Printf("lp -o number-up=2 %v # %v\n", unite(outDir, backFile, backPages(pages)), backPages(pages))
	}
	clean(outDir)
}
