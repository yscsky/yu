package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

type args struct {
	Mod, Pn string
}

func main() {
	mod := flag.String("m", "", "go mod name")
	build := flag.Bool("b", true, "build after create files")
	flag.Parse()
	if *mod == "" {
		log.Println("must set -m git url")
		return
	}
	pn, err := getProjectName()
	if err != nil {
		log.Println("getProjectName err:", err)
		return
	}
	createGitIgnore(pn)
	createMain(*mod)
	createInternal(*mod, pn)
	createGoMod(*mod)
	createDockerfile(pn)
	createReadme()
	if *build {
		goBuild()
	}
}

func getProjectName() (pn string, err error) {
	dir, err := os.Getwd()
	if err != nil {
		return
	}
	pn = filepath.Base(dir)
	return
}

func createGitIgnore(pn string) {
	createTemplateFile(".gitignore", gitignore, pn)
}

func createTemplateFile(fileName, temp string, data interface{}) {
	file, err := os.Create(fileName)
	if err != nil {
		log.Println("create", fileName, "err:", err)
		return
	}
	defer file.Close()
	if err = template.Must(template.New(fileName).Parse(temp)).Execute(file, data); err != nil {
		log.Println(fileName, "template err:", err)
		return
	}
	log.Println(fileName, "create ok")
}

func createMain(mod string) {
	createTemplateFile("main.go", mainTemp, &args{Mod: mod})
}

func createInternal(mod, pn string) {
	if err := createDir("internal/app"); err == nil {
		createTemplateFile("internal/app/app.go", appgo, &args{Mod: mod, Pn: pn})
	}
	if err := createDir("internal/db"); err == nil {
		createTemplateFile("internal/db/db.go", dbgo, &args{Mod: mod})
		createTemplateFile("internal/db/sqls.go", sqlsgo, nil)
	}
	if err := createDir("internal/model"); err == nil {
		createTemplateFile("internal/model/const.go", constgo, nil)
		createTemplateFile("internal/model/model.go", modelgo, nil)
	}
	if err := createDir("internal/config"); err == nil {
		createTemplateFile("internal/config/config.go", configgo, nil)
	}
}

func createDir(dirName string) (err error) {
	if err = os.MkdirAll(dirName, os.ModePerm); err != nil {
		log.Println(dirName, "MkdirAll err:", err)
	}
	return
}

func createGoMod(mod string) {
	if err := exec.Command("go", "mod", "init", mod).Run(); err != nil {
		log.Println("go mod init fail, err:", err)
		return
	}
	log.Println("go mod init ok")
	if err := exec.Command("go", "mod", "tidy").Run(); err != nil {
		log.Println("go mod tidy fail, err:", err)
		return
	}
	log.Println("go mod tidy ok")
}

func createDockerfile(pn string) {
	createTemplateFile("Dockerfile", dockerfile, pn)
}

func createReadme() {
	createTemplateFile("readme.md", readme, nil)
}

func goBuild() {
	if err := exec.Command("go", "build").Run(); err != nil {
		log.Println("go build fail, err:", err)
		return
	}
	log.Println("go build ok")
}
