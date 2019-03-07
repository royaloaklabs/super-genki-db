package freq

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

var DataTable map[string]float64

const (
	WordLexFile      = "Jap.Freq.2.txt"
	CorpusWordsFile  = "internet-jp-forms.num"
	CorpusLemmasFile = "internet-jp.num"
)

var mutex sync.Mutex
var Wg sync.WaitGroup

func addToDataTable(w string, f float64) {
	mutex.Lock()
	defer mutex.Unlock()

	if val, ok := DataTable[w]; ok {
		DataTable[w] = val + f
	} else {
		DataTable[w] = f
	}
}

func BuildFrequencyData() {
	fmt.Println("[INFO] Building Frequency Data")
	DataTable = make(map[string]float64)

	freqFunc := func(path string, delim string, lineSkip, wordIndex, dataIndex int, hasMultipleData bool) {
		filePath := fmt.Sprintf("%s%c%s", "data", os.PathSeparator, path)
		f, err := os.Open(filePath)
		if err != nil {
			panic(err)
		}

		scanner := bufio.NewScanner(f)

		// skip past initial data
		for i := 0; i < lineSkip; i++ {
			scanner.Scan()
		}

		// process data
		for scanner.Scan() {
			line := scanner.Text()
			sections := strings.Split(line, delim)

			word := sections[wordIndex]
			var nums []string
			if hasMultipleData {
				nums = sections[dataIndex:]
			} else {
				nums = sections[dataIndex : dataIndex+1]
			}

			var frequency float64
			for _, val := range nums {
				n, err := strconv.ParseFloat(val, 64)
				if err != nil {
					continue
				}
				frequency += n
			}

			frequency /= float64(len(nums))

			addToDataTable(word, frequency)
		}

		Wg.Done()
	}

	Wg.Add(1)
	go freqFunc(WordLexFile, "\t", 1, 0, 1, true)
	Wg.Add(1)
	go freqFunc(CorpusWordsFile, " ", 4, 2, 1, false)
	Wg.Add(1)
	go freqFunc(CorpusLemmasFile, " ", 4, 2, 1, false)
}
