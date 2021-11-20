package util

import (
	"bufio"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

type AesFileEncode struct {
	PwdKey []byte
}

func (a AesFileEncode) encode(sourceFile, destinationFile string) error {
	if IsFileExists(destinationFile) {
		return errors.New(destinationFile + " 目标文件已经存在，不能处理")
	}
	return EncryptFile(a.PwdKey, sourceFile, destinationFile)
}

func (a AesFileEncode) decode(sourceFile, destinationFile string) error {
	if IsFileExists(destinationFile) {
		return errors.New(destinationFile + " 目标文件已经存在，不能处理")
	}

	f, err := os.Open(sourceFile)
	if err != nil {
		fmt.Println("未找到文件")
		return err
	}

	br := bufio.NewReader(f)
	readString, err := br.ReadString('\n')
	f.Close()

	if err != nil {
		return err
	}

	key := a.PwdKey

	if strings.Index(readString, "0xYN7") == 0 {
		// 读取 JSON
		var headTitle FileHeadTitle
		err = json.Unmarshal([]byte(readString[6:]), &headTitle)
		if len(headTitle.PKey) > 0 {
			key = []byte(headTitle.PKey)
		}
	}

	return DecryptFile(key, sourceFile, destinationFile)
}

func padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

func unPadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("加密字符串错误！")
	}
	unPadding := int(data[length-1])
	return data[:(length - unPadding)], nil
}

func AesEncrypt(key []byte, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	encryptBytes := padding(data, blockSize)
	//初始化加密数据接收切片
	crypted := make([]byte, len(encryptBytes))
	//使用cbc加密模式
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	//执行加密
	blockMode.CryptBlocks(crypted, encryptBytes)
	return crypted, nil
}

//AesDecrypt 解密
func AesDecrypt(key []byte, data []byte) ([]byte, error) {
	//创建实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//获取块的大小
	blockSize := block.BlockSize()
	//使用cbc
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	//初始化解密数据接收切片
	crypted := make([]byte, len(data))
	//执行解密
	blockMode.CryptBlocks(crypted, data)
	//去除填充
	crypted, err = unPadding(crypted)
	if err != nil {
		return nil, err
	}
	return crypted, nil
}

//EncryptByAes Aes加密 后 base64 再加
func EncryptByAes(key []byte, data []byte) (string, error) {
	res, err := AesEncrypt(key, data)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(res), nil
}

//DecryptByAes Aes 解密
func DecryptByAes(key []byte, data string) ([]byte, error) {
	dataByte, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	return AesDecrypt(key, dataByte)
}

// EncryptFile 加密文件
// sourceFile 原文件地址
// destinationFile 目标文件地址
func EncryptFile(key []byte, sourceFile, destinationFile string) error {
	f, err := os.Open(sourceFile)
	if err != nil {
		fmt.Println("未找到文件")
		return err
	}
	defer f.Close()

	fInfo, _ := f.Stat()
	fLen := fInfo.Size()
	if fLen == 0 {
		// 文件是空的, 直接复制
		os.Create(destinationFile)
		return nil
	}
	fmt.Println("待处理文件大小:", fLen)
	maxLen := 1024 * 1024 * 50 // 50mb  每 50mb 进行加密一次
	var forNum int64 = 0
	getLen := fLen

	if fLen > int64(maxLen) {
		getLen = int64(maxLen)
		forNum = fLen / int64(maxLen)
		fmt.Println("需要加密次数：", forNum+1)
	}

	//加密后存储的文件
	ff, err := os.OpenFile(destinationFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("文件写入错误")
		return err
	}
	defer ff.Close()

	headTitle := FileHeadTitle {
		Algorithm: "aes",
		PKey:      string(key),
	}
	headTitleBytes, _ := json.Marshal(headTitle)
	title := "0xYN7 " + string(headTitleBytes) + "\n"
	buf := bufio.NewWriter(ff)
	buf.WriteString(title)
	buf.Flush()

	//循环加密，并写入文件
	for i := 0; i < int(forNum+1); i++ {
		a := make([]byte, getLen)
		n, err := f.Read(a)
		if err != nil {
			fmt.Println("文件读取错误")
			return err
		}
		getByte, err := EncryptByAes(key, a[:n])
		if err != nil {
			fmt.Println("加密错误")
			return err
		}
		// 换行处理，有点乱了，想到更好的再改
		getBytes := append([]byte(getByte), []byte("\n")...)
		//写入
		buf := bufio.NewWriter(ff)
		buf.WriteString(string(getBytes[:]))
		buf.Flush()
	}
	ffInfo, _ := ff.Stat()
	fmt.Printf("文件加密成功，生成文件名为：%s，文件大小为：%v Byte \n", ffInfo.Name(), ffInfo.Size())
	return nil
}

//DecryptFile 文件解密
func DecryptFile(key []byte, sourceFile string, destinationFile string) error {
	f, err := os.Open(sourceFile)
	if err != nil {
		fmt.Println("未找到文件")
		return err
	}
	defer f.Close()
	fInfo, _ := f.Stat()
	fmt.Println("待处理文件大小:", fInfo.Size())

	br := bufio.NewReader(f)
	ff, err := os.OpenFile(destinationFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("文件写入错误")
		return err
	}
	defer ff.Close()
	num := 0
	//逐行读取密文，进行解密，写入文件
	for {
		num = num + 1
		readline, err := br.ReadString('\n')
		if err != nil {
			break
		}
		if num ==1 && strings.Index(readline, "0xYN7") == 0 {
			continue
		}
		getByte, err := DecryptByAes(key, readline)
		if err != nil {
			fmt.Println("解密错误")
			return err
		}

		buf := bufio.NewWriter(ff)
		buf.Write(getByte)
		buf.Flush()
	}
	fmt.Println("解密次数：", num)
	ffInfo, _ := ff.Stat()
	fmt.Printf("文件解密成功，生成文件名为：%s，文件大小为：%v Byte \n", ffInfo.Name(), ffInfo.Size())
	return nil
}
