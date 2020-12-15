package main

import (
	"github.com/spf13/viper"
	"github.com/urfave/cli"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	conf string
)

func init() {
	cmd()
	for _, v := range os.Args {
		if v == "-h" || v == "--help" {
			os.Exit(1)
		}
	}
	loadYaml()
}

func cmd() {
	app := cli.NewApp()
	app.Name = "gitlab同步工具"
	app.Action = func(c *cli.Context) error {
		var showhelp error
		if c.NumFlags() == 0 {
			showhelp = cli.ShowAppHelp(c)
			os.Exit(1)
		}
		return showhelp
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "c",
			Usage:       "配置文件路径",
			Destination: &conf,
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Panicf("读取参数失败:%v", err)
	}
}

func loadYaml() {
	paths, fileName := filepath.Split(conf)
	sp := strings.Split(fileName, ".")
	if len(sp) != 2 {
		log.Panicf("文件名错误%s\n", fileName)
	}
	log.Printf("文件路径：%s\n", paths)
	log.Printf("文件名字%s\n", sp[0])
	log.Printf("文件后缀%s\n", sp[1])
	confyaml = viper.New()
	confyaml.SetConfigName(sp[0])
	confyaml.SetConfigType(sp[1])
	confyaml.AddConfigPath(paths)
	err := confyaml.ReadInConfig()
	if err != nil {
		log.Panic(err)
	}
}
