package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"unicode"

	"github.com/chai2010/gettext-go/gettext/po"
)

func main() {
	poFile, err := po.Load("../zh_CN/gccgo.po")
	if err != nil {
		log.Fatal(err)
	}
	for _, msg := range poFile.Messages {
		genTestFile(msg.MsgId)
	}
	fmt.Printf("Done\n")
}

func genTestFile(msgId string) {
	os.MkdirAll("../tests/", 0666)
	filename := "../tests/" + genFileNameFromMsgId(msgId)
	data, _ := ioutil.ReadFile(filename)
	if len(data) != 0 && string(data) != "TODO" {
		fmt.Printf("skip %v\n", filename)
		return
	}
	if err := ioutil.WriteFile(filename, []byte(`TODO`), 0666); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("generate %v\n", filename)
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
