// Copyright 2014 Golang-zh. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
测试代码生成工具

该工具需要在当前目录执行(`go run mktests.go`), 生成的测试文件放在 `../tests/` 目录.

处理的要点:

	1. 读取并解析 `../zh_CN/gccgo.po` 文件, 从中提取需要翻译的全部信息
	2. 将要翻译的 `MsgId` 映射为 `../tests/xxx.go` 文件名, 判断映射文件名是唯一的, 不唯一则停止!
	3. 映射的规则: 将删除掉开头和末尾的空白字符, 将非字母/数字映射为`_`, 具体见 `genFileNameFromMsgId` 函数
	4. 遍历 `../tests/?.go` 目录文件列表, 标记出没有被 `po.Messages` 映射到的文件名(`MsgId`可能已经变化或被删除).
	5. 依次处理`MsgId`映射后的每个唯一文件:
		a. 读取 `../tests/xxx.go` 文件
		b. 如果 文件内容为 空, 或为 "TODO", 则 生成新的测试文件
		c. 如果 文件内容以 `//po:MsgId "${MsgId}"` 开头, 则说明不需要更新, 打印信息 `skip existing ../test/xxx.go`
		d. 如果 文件内容以 `//po:MsgId` 开头, 但是 `${MsgId}` 部分不一致, 则说明文件格式被破坏, 需要打印提示信息
		e. 如果 文件内容是其他, 则 增加 `mktest.go` 生成的头部注释
	6. 退出前需要打印的信息
		0. `MsgId` 映射文件名不唯一, 停止处理, 需要重写 `mktest.go` 生成脚本
		a. 跳过的文件, 格式: skip existing ../tests/xxx.go
		b. 更新的文件, 格式: generate ../tests/xxx.go
		c. 无效的文件, 格式: invalid ../tests/xxx.go
		d. 坏掉的文件, 格式: bad ../tests/xxx.go

生成的测试文件的头部注释样例:

	//po:MsgId "ambiguous method %s%s%s"
	//po:MsgStr "有歧义的方法%s%s%s"
	//po:MsgStr ""

	package p

补充说明:

	1. PoMessage 中的其他信息需要从 po 文件同步, 不仅仅需要验证 `${MsgId}`
	2. 翻译只能在 `../tests/xxx.go` 文件进行
	3. 合并 `../tests/xxx.go` 中的 `MsgStr` 到 po 文件
*/
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
	testsDir      = "../tests/"           // 测试文件目录
	oldFileMap    = make(map[string]bool) // `testsDir` 原有的文件
	newFileMap    = make(map[string]bool) // 新生成的文件
	skipedFileMap = make(map[string]bool) // 跳过的文件
)

func main() {
	scanTestFiles()

	po, err := po.Load("../zh_CN/gccgo.po")
	if err != nil {
		log.Fatal(err)
	}
	assertPoMessageIdIsUnique(po.Messages)
	sort.Sort(byMsgId(po.Messages))

	for _, msg := range po.Messages {
		genTestFile(msg)
	}

	printSkipedFiles()
	printInvalidFiles()

	fmt.Printf("Done\n")
}

func assertPoMessageIdIsUnique(poMessages []po.Message) {
	nameMap := make(map[string][]po.Message)
	for _, msg := range poMessages {
		name := genFileNameFromMsgId(msg.MsgId)
		nameMap[name] = append(nameMap[name], msg)
	}
	if len(nameMap) != len(poMessages) {
		for name, msgList := range nameMap {
			if len(msgList) > 1 {
				fmt.Printf("Name conflicted: %s\n", name)
				for i := 0; i < len(msgList); i++ {
					fmt.Printf("\tMsgId: %s, line %d\n", msgList[i].MsgId, msgList[i].StartLine)
				}
			}
		}
		os.Exit(1)
	}
}

func scanTestFiles() {
	filenames, err := filepath.Glob(testsDir + "*.go")
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < len(filenames); i++ {
		oldFileMap[filepath.Base(filenames[i])] = true
	}

	// do we need skip this file?
	for i := 0; i < len(filenames); i++ {
		d, _ := ioutil.ReadFile(filenames[i])
		if bytes.HasPrefix(d, []byte("//po:MsgId")) {
			skipedFileMap[filepath.Base(filenames[i])] = true
			touchTestFile(filepath.Base(filenames[i]))
		}
	}
}

func printSkipedFiles() {
	for filename, _ := range skipedFileMap {
		fmt.Printf("skip existing %s%s\n", testsDir, filename)
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

func isTestFileExists(filename string) bool {
	_, ok := skipedFileMap[filename]
	return ok
}

func touchTestFile(filename string) {
	newFileMap[filename] = true
}

// 生成测试文件
func genTestFile(msg po.Message) {
	if isTestFileExists(genFileNameFromMsgId(msg.MsgId)) {
		return
	}

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

// 根据 MsgId 生成 文件名(假设生成的文件名是唯一的)
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

// 更新头部翻译注释
func updateHeaderComments(msg po.Message, data []byte) (newData []byte) {
	defer func() {
		for bytes.Contains(newData, []byte("\n\n\n")) {
			newData = bytes.Replace(newData, []byte("\n\n\n"), []byte("\n\n"), -1)
		}
	}()
	headerComments := genHeaderComments(msg)

	// empty file
	if d := string(data); d == "" || d == "TODO" {
		newData = append(headerComments, []byte("package p")...)
		return
	}
	// user comments
	if !bytes.HasPrefix(data, []byte("//po:MsgId")) {
		newData = append(headerComments, data...)
		return
	}
	// MsgId changed, why?
	if bytes.HasPrefix(data, []byte("//po:MsgId")) && !bytes.HasPrefix(data, []byte(genHeaderFirstLine(msg))) {
		newData = append(headerComments, data...)
		return
	}
	// rewrite PoMessage comments
	if idxDoc := bytes.Index(data, []byte("\n\n")); idxDoc >= 0 {
		newData = append(headerComments, data[idxDoc:]...)
		return
	}
	newData = append(headerComments, data...)
	return
}

/*
生成的头部注释

头部注释含全部的翻译信息.

样例:

	//po:MsgId "value computed is not used"
	//po:MsgStr ""
*/
func genHeaderComments(msg po.Message) []byte {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%s\n", genHeaderFirstLine(msg))

	// 切割翻译后的字符串(可能多行)
	msgStr := msg.MsgStr
	msgStr = strings.Replace(msgStr, `\n`, "\n", -1)
	msgStr = strings.TrimSpace(msgStr)
	msgStrLines := strings.Split(msgStr, "\n")

	// 删除尾部的空行
	if n := len(msgStrLines); n > 0 && msgStrLines[n-1] == "" {
		msgStrLines = msgStrLines[:n-1]
	}

	// 如果没有数据则, 新建一个空翻译
	if len(msgStrLines) == 0 {
		msgStrLines = []string{""}
	}

	// 输出翻译部分
	for i := 0; i < len(msgStrLines); i++ {
		fmt.Fprintf(&buf, "//po:MsgStr \"%s\"\n", msgStrLines[i])
	}
	fmt.Fprintf(&buf, "\n")
	return buf.Bytes()
}

// 生成文件的第一行信息
//
// 格式:
//	//po:MsgId "ambiguous method %s%s%s"
func genHeaderFirstLine(msg po.Message) string {
	return fmt.Sprintf("//po:MsgId \"%s\"", msg.MsgId)
}

// 用于 []po.Message 排序
type byMsgId []po.Message

func (d byMsgId) Len() int           { return len(d) }
func (d byMsgId) Less(i, j int) bool { return d[i].MsgId < d[j].MsgId }
func (d byMsgId) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
