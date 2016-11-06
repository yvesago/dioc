package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"testing"
)

func createTestFile() {

	lines := []string{
		"line 1\n",
		"line 2\n",
		"line 3\n",
		"line 4\n",
		"line 5\n",
	}

	f, _ := os.Create("_test.log")
	w := bufio.NewWriter(f)

	for _, l := range lines {
		w.WriteString(l)
	}

	w.Flush()
}

func removeTestFile() {
	var err = os.Remove("_test.log")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}

}

func _TestSyncList(t *testing.T) {
	Surveys = map[string]*regexp.Regexp{"test1": nil, "test2": nil}
	Surveys["test1"] = nil
	Surveys["test2"] = nil
	fmt.Println(Surveys)
	currentList := []string{"test 2", "test 1", "test 3"}
	fmt.Println(currentList)
//	SyncList(currentList)
	fmt.Println(Surveys)
}

func TestFullSearch(t *testing.T) {
	defer removeTestFile()
	Surveys["test1"] = regexp.MustCompile("line 3")

	createTestFile()

	lines := FullSearch("_test.log")
	fmt.Println(lines)
	for _, l := range lines {
		fmt.Println(l.line)
	}

}

func TestGetLines(t *testing.T) {
	defer removeTestFile()

	createTestFile()

	lines := GetLines("_test.log")
	fmt.Println(lines)
	for _, l := range lines {
		fmt.Println(l)
	}

}
