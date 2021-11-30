package main

import (
	"AesFileUtil/util"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"path"
	"strings"
)

func main() {
	app := &cli.App{
		Name:      "fde",
		Usage:     "file decode and encode",
		UsageText: "fde [global options] [command options] [arguments...]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "key",
				Aliases: []string{"k"},
				Usage:   "加解密用的 key",
				Value:   "@shuangguidaidan",
			},
			&cli.BoolFlag{
				Name:  "f",
				Usage: "如果目标文件已经存在，直接删除，并创建新文件",
			}, &cli.BoolFlag{
				Name:  "d",
				Usage: "标识此次操作为解码，默认为编码",
			}, &cli.BoolFlag{
				Name:  "r",
				Usage: "如果源文件地址是一个目录，则将该目录以及子目录下的所有文件进行编码或解密",
			}, &cli.BoolFlag{
				Name:  "i",
				Usage: "忽略错误，继续执行",
			},
		},
		Action: func(c *cli.Context) error {
			if c.NArg() == 2 {
				isDecode := c.Bool("d")
				privateKey := c.String("key")
				isRecursion := c.Bool("r")
				sourceFile := c.Args().Get(0)
				destinationFile := c.Args().Get(1)

				aesFileEncode := util.AesFileEncode{PwdKey: []byte(privateKey)}

				if !isRecursion {
					if isDecode {
						_ = aesFileEncode.Decode(sourceFile, destinationFile)
					} else {
						_ = aesFileEncode.Encode(sourceFile, destinationFile)
					}
				}
				// 递归
				if !(util.IsDir(sourceFile) && util.IsDir(destinationFile)) {
					log.Println("递归模式两个路径必须都是文件夹")
				}

				files := make([]string, 0)

				treedir(sourceFile, &files)

				for _, file := range files {
					if isDecode {

						err := aesFileEncode.Decode(file, path.Join(destinationFile, strings.Replace(file, sourceFile, "", 1)))
						if err != nil {
							log.Println(err)
						}
					} else {
						realDestinationFile := path.Join(destinationFile, strings.Replace(file, sourceFile, "", 1))
						if util.IsDir(file) && !util.IsFileExists(realDestinationFile) {
							os.MkdirAll(realDestinationFile, os.ModePerm)
							continue
						}
						err := aesFileEncode.Encode(file, realDestinationFile)
						if err != nil {
							log.Println(err)
						}
					}
				}

			} else {
				log.Println("参数数量不对，只能有 2 个参数")
			}
			return nil
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func treedir(fpath string, files *[]string) {
	// 获取fileinfo
	if finfo, err := os.Stat(fpath); err == nil {
		// 判断是不是目录 如果不是目录而是文件  打印文件path并跳出递归
		if !finfo.IsDir() {
			return
		} else {
			// 是目录的情况 打印目录path
			f, _ := os.Open(fpath) // 通过目录path open一个file
			defer f.Close()
			names, _ := f.Readdirnames(0) // 通过file的Readdirnames 拿到当前目录下的所有filename
			for _, name := range names {
				if strings.Index(name, ".") == 0 {
					continue
				}

				newpath := path.Join(fpath, name) // 遍历names 拼接新的fpath
				*files = append(*files, newpath)
				treedir(newpath, files) // 递归
			}
		}
	}
}
