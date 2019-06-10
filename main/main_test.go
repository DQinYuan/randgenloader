package main

import (
	"github.com/DQinYuan/randgenloader"
	"testing"
)

func TestStartServer(t *testing.T) {
	randgenloader.ConfPath = "/home/dqyuan/language/Go/projects/randgenloader"
	randgenloader.RmPath = "/home/dqyuan/language/Mysql/randgenx"
	randgenloader.ResultPath = "/home/dqyuan/language/Go/projects/randgenloader"
	main()
}
