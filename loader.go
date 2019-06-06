package randgenloader

import (
	"bytes"
	"fmt"
	"github.com/pmezard/go-difflib/difflib"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type RandgenLoader struct {
	testName string
	// 存放yyzz文件的path
	confPath string
	//  存放结构的path
	resultPath string
	//  randgen主目录
	rmPath string

	cachedQueries []string
}

const confPath = "/root/conf"
const rmPath = "/root/randgenx"
const resultPath = "/root/result"

func (rl *RandgenLoader) Init(testName string) {
	rl.confPath = confPath
	rl.resultPath = resultPath
	rl.rmPath = rmPath
	rl.testName = testName
}

func (rl *RandgenLoader) LoadData(zzContent string, yyContent string) (sqls []string, err error) {
	zzPath := fmt.Sprintf(filepath.Join(rl.confPath, "%s.zz"), rl.testName)
	yyPath := fmt.Sprintf(filepath.Join(rl.confPath, "%s.yy"), rl.testName)
	ioutil.WriteFile(zzPath, []byte(zzContent), os.ModePerm)
	ioutil.WriteFile(yyPath, []byte(yyContent), os.ModePerm)

	rPath := filepath.Join(rl.resultPath, rl.testName)
	_, err = execShell(rl.rmPath, "perl", "gentest.pl",
		fmt.Sprintf("--dsn=dummy:file:%s", rPath),
		fmt.Sprintf("--gendata=%s", zzPath),
		fmt.Sprintf("--grammar=%s", yyPath))

	if err != nil {
		return nil, err
	}

	f, err := os.Open(rPath)
	if err != nil {
		return nil, err
	}
	sqlBytes, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	data, grammar := splitToDataAndGrammar(sqlBytes)
	rl.cachedQueries = grammar

	return data, nil
}

func splitToDataAndGrammar(totalContent []byte) (data []string, grammar []string) {
	content := string(totalContent)

	gendataAndGrammar := strings.Split(content, "/* follow is grammar sql */;\n")

	return strings.Split(gendataAndGrammar[0], "\n"), strings.Split(gendataAndGrammar[1], "\n")
}

func (rl *RandgenLoader) Query() (sqls []string) {
	return rl.cachedQueries
}

// r1为Mysql输出  r2为TiDB输出
func (rl *RandgenLoader) Compare(r1 string, r2 string) (comment string, consistent bool) {
	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(r1),
		B:        difflib.SplitLines(r2),
		FromFile: "Mysql",
		ToFile:   "Tidb",
	}
	text, _ := difflib.GetUnifiedDiffString(diff)

	return text, text == ""
}

//执行shell命令
func execShell(dir string, s string, args ...string) (string, error) {
	cmd := exec.Command(s, args...)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	cmd.Dir = dir
	err := cmd.Run()

	return out.String(), err
}
