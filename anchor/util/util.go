// Copyright 2018 The Liman Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package util

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/golang/glog"
	"golang.org/x/crypto/bcrypt"
)

// HTTPGet executes get http method and return response
func HTTPGet(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		glog.Errorf("URL=%s; Err=%s", url, err)
		return []byte{}, err
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

//GeneratePassword for apiKey and cookieValue
func GeneratePassword(l int) string {
	rand.Seed(time.Now().UnixNano())
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

	b := make([]rune, l)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

//Version getting git commit id
func Version() (string, error) {
	var version string

	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			version = scanner.Text()
		}
	}()

	err = cmd.Start()
	if err != nil {
		return "", err
	}

	err = cmd.Wait()
	if err != nil {
		return "", err
	}

	return version, nil
}

//HashPasswordAndSave Hashing root password and storing them.
func HashPasswordAndSave(p string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(p), 14)
	path, _ := os.Getwd()
	if err != nil {
		return "", err
	}

	err = os.Mkdir(path+"/data", 0777)
	if err != nil {
		glog.Error("/data folder already exist. Skipping.")
	}

	err = ioutil.WriteFile(path+"/data/pass", b, 0644)
	if err != nil {
		return "", err
	}

	return string(b), nil

}

//ReadPassword read hashed password from file
func ReadPassword() string {
	path, _ := os.Getwd()
	h, err := ioutil.ReadFile(path + "/data/pass")
	if err != nil {
		return ""
	}

	return string(h)
}

//CheckPass matchs hash value with pass
func CheckPass(p, h string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(h), []byte(p))
	return err == nil
}

// StringAdd add two number string and return string
func StringAdd(num1, num2 string) (string, error) {

	var err error
	var num1f, num2f float64

	if num1 == "" {
		num1f = 0
	} else {

		num1f, err = strconv.ParseFloat(num1, 32)
		if err != nil {
			return "", err
		}
	}

	if num2 == "" {
		num2f = 0
	} else {
		num2f, err = strconv.ParseFloat(num2, 32)
		if err != nil {
			return "", err
		}
	}
	return fmt.Sprintf("%f", num1f+num2f), nil
}
