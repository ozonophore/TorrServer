package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func main() {
	dir, _ := os.Getwd()
	if _, err := os.Stat("web/build/static"); os.IsNotExist(err) {
		os.Chdir("web")
		if run("yarn") != nil {
			os.Exit(1)
		}
		if run("yarn", "run", "build") != nil {
			os.Exit(1)
		}
		os.Chdir(dir)
	}

	compileHtml := "web/build/"
	srcGo := "server/web/pages/"

	run("rm", "-rf", srcGo+"template/pages")
	run("cp", "-r", compileHtml, srcGo+"template/pages")
}

func run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func cleanName(fn string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return strings.Title(reg.ReplaceAllString(fn, ""))
}
