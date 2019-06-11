package randgenloader

import (
	"fmt"
	"github.com/pmezard/go-difflib/difflib"
	"io/ioutil"
	"log"
	"path/filepath"
	"testing"
)

func TestWriteFile(t *testing.T) {
	t.SkipNow()
	ioutil.WriteFile("test.tt", []byte("opopii"), 0777)
}

const (
	text1 = "Lorem\nkkkk"
	text2 = "Lorem\nkkkk"
)

func TestDiff(t *testing.T) {
	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(text1),
		B:        difflib.SplitLines(text2),
		FromFile: "Mysql",
		ToFile:   "Tidb",
	}
	text, _ := difflib.GetUnifiedDiffString(diff)
	fmt.Println(text == "")
}

func TestSplitDataAndGrammar(t *testing.T) {
	allContent, _ := ioutil.ReadFile("./haha")

	data, grammar := splitToDataAndGrammar(allContent)
	fmt.Println(data)   // [ddawww ooooo ]
	fmt.Println(len(data))  // 3
	fmt.Println(grammar)   // [amibgy komaisa ]
	fmt.Println(len(grammar))  // 3
}

func pwd() string {
	s, _ := filepath.Abs(".")
	return s
}

func TestLoader(t *testing.T) {
	rl := new(RandgenLoader)

	rl.Init("haha")
	rl.ConfPath = pwd()
	rl.ResultPath = pwd()
	rl.RmPath = "/home/dqyuan/language/Mysql/randgenx"

	zzContent, _ := ioutil.ReadFile("./yyzzs/example.zz")
	yyContent, _ := ioutil.ReadFile("./yyzzs/example.yy")
	ddls, err := rl.LoadData(string(zzContent), string(yyContent), "testtest", "10")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(ddls)
	fmt.Printf("ddl len: %d\n", len(ddls))
	fmt.Println("====================")
	dqls := rl.Query()
	fmt.Println(dqls)
	fmt.Printf("dql len: %d\n", len(dqls))

	counter := 0
	for _, c := range dqls {
		if c == "" {
			counter++
		}
	}

	fmt.Println(counter)
}
