package main

import (
	"unicode/utf8"
	//"bufio"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	//"reflect"
	"strings"

	"github.com/axgle/mahonia"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	sourceFile = kingpin.Flag("src", "Use -s <源 IP地址>").Default("").Short('s').String()
	//targetFile = kingpin.Flag("tgt", "Use -t <目标IP地址>").Default("").Short('t').String()
)

//主体程序
func main() {

	//解析命令参数
	kingpin.Parse()

	filename := *sourceFile
	if filename == "" {
		fmt.Println("即将在当前目录查找 txt 文本,并转出excel")
		dir := getCurrentFilePath()
		fmt.Print(dir)
		readDir(dir)
		fmt.Println("complete!  press any key to close this window!")
		var in string
		fmt.Scanf("%s", &in)
	} else {
		fmt.Println("即将对", "XXX", ",并转出excel")
		str := readFile(filename)
		if str == "" {
			fmt.Println("readfile error!!!!!!!")
			return
		}
		txt2Excel(filename, str)
	}
}

//获取当前文件路径
func getCurrentFilePath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	dir, _ := filepath.Split(path)
	return dir
}

//传递目录遍历出 txt 文本文件,并转出为excel。
func readDir(dir string) {
	rd, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println("read dir error!")
		return
	}
	for _, file := range rd {
		filename := file.Name()
		if strings.Contains(filename, ".txt") {
			fmt.Println("filename: ", filename)
			str := readFile(filename)
			if str == "" {
				fmt.Println("readfile error!!!!!!!")
				return
			}
			txt2Excel(filename, str)
		}
	}
}

//txt转excel
func txt2Excel(filename string, str string) {
	xlsname := strings.Trim(filename, ".txt") + ".xls"
	f, err := os.Create(xlsname)
	if err != nil {
		fmt.Println("create " + xlsname + "failed!!!!!")
	}
	defer f.Close()
	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
	w := csv.NewWriter(f)

	sli := getRowsSli(str) //sli:[  12  45 8,   456 8  9,...]
	for _, v := range sli {
		val := formatString(v) //val: [12 45 8]
		w.Write(val)
	}
	w.Flush()
}

//txt读取，为gbk编码的自动转utf8
func readFile(filename string) string {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("there is no source.txt")
		return ""
	}
	f := string(file)
	isutf8 := utf8.ValidString(f)
	if isutf8 == false {
		utf8 := mahonia.NewDecoder("gbk").ConvertString(f)
		return utf8
	} else {
		return f
	}
}

//txt中换行是“\n”
func getRowsSli(str string) []string {
	sli := strings.Split(str, "\n")
	return sli
}

// 以1个或多个空白字符分隔字符串s并返回slice
func formatString(str string) []string {
	tmp := ""
	for _, s := range str {
		tmp = tmp + string(s)
	}
	rtnSli := strings.Fields(tmp)
	return rtnSli
}
