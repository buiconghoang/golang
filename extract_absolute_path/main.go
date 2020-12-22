package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {
	// testContainsLinuxPath()
	// testsplitApostropheInLinuxPath()
	// testnormalizeSlashWindowPath()
	// testsaveExtractPath()
	testsaveExtractPath()
}

func testContainsLinuxPath() {
	command := "mv /home/user/a.py /home/user2/"
	fmt.Println(isAbleToContainsALinuxPath(command))
}

func testsplitApostropheInLinuxPath() {
	command := `cp /home/user/s\'sba/oke/a.py 's\ ab/su\'bspace/abc'`
	fmt.Printf("%q\n", splitApostropheAndSpaceInLinuxPath(command))
}

func testnormalizeSlashWindowPath() {
	command := "c:\\\\\\\\\\abc\\def//hihi//hha"
	fmt.Println(normalizeSlashWindowPath(command))
}

func testsaveExtractPath() {
	datapath := "linux_paths.txt"
	saveExtractPath(datapath, "extract_linux_path.csv", extractPath)
}

func isAbleToContainsALinuxPath(command string) bool {
	re := regexp.MustCompile(`(^/|\s\/)`)
	result := re.FindAllString(command, -1)
	fmt.Println(result)
	if len(result) > 0 {
		return true
	}
	return false
}

func splitApostropheAndSpaceInLinuxPath(command string) []string {
	// go does not support negative look-behind
	// => use findIndex to find index of slash and apostrophe
	// then slice on command
	// => stop-1 mean remove  slash or apostrophe was found

	re := regexp.MustCompile(`[^\\][\'\s]+`)
	result := []string{}
	commandLen := len(command)
	partIdxes := re.FindAllIndex([]byte(command), -1)
	fmt.Println(partIdxes)
	start := 0
	stop := 0
	for _, idx := range partIdxes {
		stop = idx[1]
		result = append(result, strings.TrimSpace(command[start:stop-1]))
		start = stop
	}
	if stop < commandLen {
		result = append(result, strings.TrimSpace(command[start:commandLen]))
	}
	return result
}

func normalizeSlashWindowPath(filePath string) string {
	re := regexp.MustCompile(`(\\+|\/+)`)
	return re.ReplaceAllString(filePath, "\\")
}

func extractPath(command string) []string {
	return []string{"path11", "path2", "path3"}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func createWriter(savePath string) *bufio.Writer {
	f, err := os.OpenFile(savePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	check(err)
	return bufio.NewWriter(f)
}

func saveExtractPath(dataPath string, savePath string, extractPath func(string) []string) {
	file, err := os.Open(dataPath)
	check(err)
	writer := createWriter(savePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		command := scanner.Text()
		writer.WriteString(";")

		writer.WriteString(command + ";")
		for _, filePath := range extractPath(command) {
			writer.WriteString(filePath + ";")
		}
		writer.WriteString("\n")
		writer.Flush()
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
