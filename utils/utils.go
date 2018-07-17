package utils

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"bytes"

	"errors"

	"io/ioutil"

	"github.com/ctcpip/notifize"
)

// checking if a directory or file exists
func PathExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// creating a new directory
func Mkdir(path string) error {
	return os.MkdirAll(path, 0755)
}

// creating a new file
func Mkfile(name string) error {
	_, err := os.Create(name)
	return err
}

// removing a file
func Rmfile(path string) error {
	return os.Remove(path)
}

// reading a file
func ReadFile(path string) (out string, err error) {
	body, err := ioutil.ReadFile(path)

	return strings.TrimSpace(string(body)), err
}

// getting correct prefix depending on the user's OS
// to get a path where a temp DB will be stored
func GetPrefix() string {
	if runtime.GOOS == "darvin" {
		return "/tmp/"
	}
	return "/dev/shm/"
}

// opening a URL in the browser
func OpenURL(url string) {
	if !strings.HasPrefix(url, "http://") &&
		!strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}

	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()

	case "darwin":
		err = exec.Command("open", url).Start()
	}

	if err != nil {
		fmt.Printf("failed to open url %s: %s\n", url, err)
	}
}

func Notify(summary, body string) {
	notifize.Display(summary, body, false, "")
}

func IsIntalled(name string) (bool, error) {
	out, err := exec.Command("whereis", name).Output()
	if err != nil {
		return false, err
	}

	data := strings.Split(string(out), ":")
	return len(data[1]) > 1, nil
}

func ShowMenu(name string, passwords string) (string, error) {
	var stdout bytes.Buffer
	stdin := strings.NewReader(passwords)

	var cmd *exec.Cmd
	if name == "dmenu" {
		cmd = exec.Command("dmenu")
	} else if name == "rofi" {
		cmd = exec.Command("rofi", "-dmenu", "-p", "select password")
	} else {
		e := errors.New("Invalid command: " + name)
		return "", e
	}

	cmd.Stdin = stdin
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(stdout.String()), nil
}

// reading a line from stdin
func ReadLine() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(line), nil
}

// generating a random name for temp DB files
func GenerateName() string {
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	name := make([]byte, 10)
	rand.Seed(time.Now().UnixNano())
	for i := range name {
		name[i] = chars[rand.Intn(len(chars))]
	}

	return string(name)
}
