package main

import (
	"flag"
	"fmt"
	"github.com/anthony-dong/template-generator/file"
	"github.com/anthony-dong/template-generator/scrpit"
	"os"
	"path/filepath"
)

var (
	gitRemoteAddress string
	cloneDir         string
	modelName        string
	helps            bool
)

func init() {
	flag.StringVar(&modelName, "mod", "", "go mod name, eg:-mod=ebike-city-report")
	flag.StringVar(&gitRemoteAddress, "git", "", "git branch, eg:-git=git@gitee.com:Anthony-Dong/template.git")
	flag.StringVar(&cloneDir, "dir", "./", "go , eg:-dir=/data/temp")
	flag.BoolVar(&helps, "h", false, "this help")
}

func main() {
	initFlag2()

	if cloneDir == "./" {
		fmt.Println("clone file in current dir")
	}
	dir, e := filepath.Abs(cloneDir)
	if e != nil {
		panic(e)
	}
	fmt.Printf("git clone %s to %s\n", gitRemoteAddress, dir)
	// clone
	scrpit.Git(gitRemoteAddress, dir)
	fmt.Printf("git clone %s success", gitRemoteAddress)
	// rename
	e = file.NewTemplate(dir, modelName)
	if e != nil {
		panic(e)
	}
	fmt.Println("build template success")

	// delete
	scrpit.Delete(join(dir, ".git"))

	// copy 配置文件
	scrpit.Copy(join(dir, "config/env.dev.ini"), join(dir, "config/env.ini"))

	fmt.Printf("new %s project success\n", modelName)
}

func join(dir, file string) string {
	return fmt.Sprintf("%s%s%s", dir, string(os.PathSeparator), file)
}

func initFlag2() {
	flag.Parse()
	if helps {
		printHelp()
		os.Exit(-1)
	}
	if modelName == "" {
		fmt.Println("please set -mod=your_project_name")
		printHelp()
		os.Exit(-1)
	}
}

func printHelp() {
	fmt.Println(`================go-build  help=====================
boot -dir=/data/temp -mod=city-demo -git=git@gitee.com:Anthony-Dong/template.git
Option:`)
	flag.PrintDefaults()
}
