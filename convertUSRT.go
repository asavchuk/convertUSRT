//This programm converts U-tube auto-generated transcriptions (.srt files),
//so that transcription will be shown in player like a one-liner solution instead of two lines.
//This can be especially usefull for big monitors. This also removes U-tube `feature`
//when auto-generated transcription shows in not straightforward order.
//Usage: place `convertUSRT.exe` in the same folder where .srt files located. Run it once.
//A processed files will be located in the `converted` folder.

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

const LineBreak = "\r\n"
const createdDir = "converted"

func main() {
	var fileName string

	fmt.Print("")
	dirname := "." + string(filepath.Separator)

	d, err := os.Open(dirname)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer d.Close()

	files, err := d.Readdir(-1)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Reading current dir: ")
	fmt.Println(CurrentFolder())

	CreateDirIfNotExist(createdDir)

	for _, file := range files {
		if file.Mode().IsRegular() {
			if filepath.Ext(file.Name()) == ".srt" {
				fileName = file.Name()
				fmt.Println("processing " + fileName)

				timeMarkers := []string{}
				transcript := []string{}
				wg := new(sync.WaitGroup)
				wg.Add(2)
				go timeMarkersSearch(fileName, &timeMarkers, wg)
				go transcriptPairsConnection(fileName, &transcript, wg)
				wg.Wait()
				writeSRT(fileName, timeMarkers, transcript)
			}
		}
	}
}

//Detecting a current folder. Returns full path.
// `dir` is the directory of the currently running file.
func CurrentFolder() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

//Create a directory if it does not exist. Otherwise do nothing.
func CreateDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
}

//Makes final processed .srt file in the "converted" folder
//which is created in the same dir from where the programm started
func writeSRT(fileName string, timeMarkers []string, transcript []string) {
	f, err := os.OpenFile("."+string(filepath.Separator)+createdDir+string(filepath.Separator)+fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	for i, marker := range timeMarkers {
		if _, err := f.WriteString(strconv.Itoa(i+1) + LineBreak + marker + LineBreak); err != nil {
			log.Println(err)
		}
		if _, err := f.WriteString(transcript[i] + LineBreak + LineBreak); err != nil {
			log.Println(err)
		}
	}
}

func transcriptPairsConnection(fileName string, transcript *[]string, wg *sync.WaitGroup) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var s string
	scanner := bufio.NewScanner(file)

	var line int
	var oneLiner string

	for scanner.Scan() {
		line++
		s = scanner.Text()
		if line == 3 {
			oneLiner += s
		}
		if line == 7 {
			oneLiner += " "
			oneLiner += s
			break
		}
	}

	*transcript = append(*transcript, oneLiner)

	line = 0
	oneLiner = ""

	for scanner.Scan() {
		line++
		s = scanner.Text()

		if line == 4 {
			oneLiner += s
		}

		if line == 8 {
			oneLiner += " "
			oneLiner += s
			*transcript = append(*transcript, oneLiner)
			line = 0
			oneLiner = ""
		}

	}

	if len(oneLiner) > 0 {
		*transcript = append(*transcript, oneLiner)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	wg.Done()
}

func timeMarkersSearch(fileName string, timeMarkers *[]string, wg *sync.WaitGroup) {

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var s string
	scanner := bufio.NewScanner(file)

	var line int

	for scanner.Scan() {
		line++
		s = scanner.Text()
		if line == 2 {
			*timeMarkers = append(*timeMarkers, s)
			break
		}
	}

	line = 0

	for scanner.Scan() {
		line++
		s = scanner.Text()

		if line == 8 {
			*timeMarkers = append(*timeMarkers, s)
			line = 0
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	wg.Done()

}
