package main

import (
	"AesFileUtil/util"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	app := &cli.App{
		Name:  "fde",
		Usage: "file decode and encode",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "key",
				Aliases: []string{"k"},
				Usage:   "加解密用的 key",
				Value:   "1234567812345678",
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
				sourceFile := c.Args().Get(0)
				destinationFile := c.Args().Get(1)

				aesFileEncode := util.AesFileEncode{PwdKey: []byte(privateKey)}

				if isDecode {
					_ = aesFileEncode.Decode(sourceFile, destinationFile)
				} else {
					_ = aesFileEncode.Encode(sourceFile, destinationFile)
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
