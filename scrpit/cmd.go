package scrpit

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

// func main() {
//	exec.Command("git", "clone", "git@gitee.com:Anthony-Dong/template.git")
//}
func Git(str, dir string) {
	gitCmd := fmt.Sprintf("git clone %s  %s", str, dir)
	Cmd(gitCmd)
}

func Run(shell string) {
	gitCmd := fmt.Sprintf("%s", shell)
	Cmd(gitCmd)
}

func Copy(src, dest string) {
	gitCmd := fmt.Sprintf("cp %s %s", src, dest)
	Cmd(gitCmd)
	fmt.Printf("copy %s to %s success\n", src, dest)
}

// delete file
func Delete(file string) {
	gitCmd := fmt.Sprintf("rm -rf %s", file)
	Cmd(gitCmd)
	fmt.Printf("delete %s file success\n", file)
}

func Cmd(cmd string) {
	command := exec.Command("/bin/bash", "-c", cmd)
	command.Stdout = os.Stdout
	err := command.Run()
	if err != nil {
		log.Fatal(err)
	}
}
