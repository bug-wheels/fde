package test

import (
	"AesFileUtil/util"
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestFileAppend(t *testing.T) {
	key := []byte("1234567812345678")

	ff, err := os.OpenFile("../hh.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		t.Log("文件写入错误")
	}
	defer ff.Close()

	// 提前写入一次，将当前的一些辅助加密的信息写入

	headTitle := util.FileHeadTitle {
		Algorithm: "aes",
		PKey:      string(key),
	}

	headTitleBytes, _ := json.Marshal(headTitle)


	title := "0xYN7 " + string(headTitleBytes) + "\n"
	buf := bufio.NewWriter(ff)
	buf.WriteString(title)
	buf.Flush()

	buf = bufio.NewWriter(ff)
	buf.WriteString("fasdfefasdf")
	buf.Flush()
}

func TestFileRead(t *testing.T) {
	f, err := os.Open("../hh.txt")
	if err != nil {
		fmt.Println("未找到文件")
		t.Error(err)
		return
	}

	br := bufio.NewReader(f)
	readString, err := br.ReadString('\n')
	f.Close()

	if err != nil {
		t.Error(err)
		return
	}

	if strings.Index(readString, "0xYN7") == 0 {
		// 读取 JSON
		t.Log(readString[6:])
		var headTitle util.FileHeadTitle
		err = json.Unmarshal([]byte(readString[6:]), &headTitle)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(headTitle)
		if len(headTitle.PKey) > 0 {
			t.Log(headTitle.PKey)
		}
	}
}