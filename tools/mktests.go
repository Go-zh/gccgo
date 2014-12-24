// Copyright 2014 Golang-zh. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"

	"github.com/chai2010/gettext-go/gettext/po"
)

var (
	testsDir   = "../tests/"
	oldFileMap = make(map[string]bool)
	newFileMap = make(map[string]bool)
)

type byMsgId []po.Message

func (d byMsgId) Len() int           { return len(d) }
func (d byMsgId) Less(i, j int) bool { return d[i].MsgId < d[j].MsgId }
func (d byMsgId) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }

func main() {
	scanTestFiles()

	po, err := po.Load("../zh_CN/gccgo.po")
	if err != nil {
		log.Fatal(err)
	}

	sort.Sort(byMsgId(po.Messages))
	for _, msg := range po.Messages {
		genTestFile(msg)
	}
	printInvalidFiles()

	fmt.Printf("Done\n")
}

func scanTestFiles() {
	filenames, err := filepath.Glob(testsDir + "*.go")
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < len(filenames); i++ {
		oldFileMap[filepath.Base(filenames[i])] = true
	}
}

func printInvalidFiles() {
	var invalidFiles []string
	for filename, _ := range oldFileMap {
		if _, ok := newFileMap[filename]; !ok {
			invalidFiles = append(invalidFiles, filename)
		}
	}
	sort.Strings(invalidFiles)
	for i := 0; i < len(invalidFiles); i++ {
		fmt.Printf("invalid %s%s\n", testsDir, invalidFiles[i])
	}
}

func touchTestFile(filename string) {
	newFileMap[filename] = true
}

func genTestFile(msg po.Message) {
	os.MkdirAll(testsDir, 0666)
	filename := testsDir + genFileNameFromMsgId(msg.MsgId)
	data, _ := ioutil.ReadFile(filename)
	data = updateHeaderComments(msg, data)
	if err := ioutil.WriteFile(filename, data, 0666); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("generate %v\n", filename)
	touchTestFile(filepath.Base(filename))
}

func genFileNameFromMsgId(msgId string) string {
	runes := []rune(msgId)
	for i := 0; i < len(runes); i++ {
		if unicode.IsLetter(runes[i]) {
			continue
		}
		if unicode.IsDigit(runes[i]) {
			continue
		}
		runes[i] = ' '
	}

	name := string(runes)
	name = strings.TrimSpace(name)
	for strings.Contains(name, "  ") {
		name = strings.Replace(name, "  ", " ", -1)
	}

	name = strings.Replace(name, " ", "_", -1)
	return name + ".go"
}

func updateHeaderComments(msg po.Message, data []byte) (newData []byte) {
	defer func() {
		for bytes.Contains(newData, []byte("\n\n\n")) {
			newData = bytes.Replace(newData, []byte("\n\n\n"), []byte("\n\n"), -1)
		}
	}()
	headerComments := genHeaderComments(msg)
	if d := string(data); d == "" || d == "TODO" {
		newData = headerComments
		return
	}
	if idxDoc := bytes.Index(data, []byte("\n\n")); idxDoc >= 0 {
		newData = append(headerComments, data[idxDoc:]...)
		return
	}
	newData = append(headerComments, data...)
	return
}

/*
Generate file header comments.

Example:
	//po:MsgId "value computed is not used"
	//
	//PoMessage {
	//	#: go/gofrontend/types.cc:7196
	//	#, c-format
	//	msgid "ambiguous method %s%s%s"
	//	msgstr "有歧义的方法%s%s%s"
	//}
*/
func genHeaderComments(msg po.Message) []byte {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "//po:MsgId \"%s\"\n", msg.MsgId)
	fmt.Fprintf(&buf, "//\n")
	fmt.Fprintf(&buf, "//PoMessage {\n")
	msgLines := strings.Split(msg.String(), "\n")
	for i := 0; i < len(msgLines); i++ {
		if msgLines[i] != "" {
			fmt.Fprintf(&buf, "//\t%s\n", msgLines[i])
		}
	}
	fmt.Fprintf(&buf, "//}\n")
	fmt.Fprintf(&buf, "\n")

	return buf.Bytes()
}
