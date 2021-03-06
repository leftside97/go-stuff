package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"bufio"
)


func getKey() string {
	keyreader := bufio.NewReader(os.Stdin)
	fmt.Println("key:")
	key, _ := keyreader.ReadString('\n')
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}


func encrypt(data []byte, keyhash string) []byte {
	block, _ := aes.NewCipher([]byte(keyhash))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Println(err)
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Println(err)
	}
	dataout := gcm.Seal(nonce, nonce, data, nil)
	return dataout
}

func decrypt(data []byte, keyhash string) []byte {
	block, err := aes.NewCipher([]byte(keyhash))
	if err != nil {
		fmt.Println(err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Println(err)
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	dataout, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		fmt.Println(err)
	}
	return dataout
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("usage: "+os.Args[0]+" [-enc|-dec] <file>")
		os.Exit(3)
	}
	switch os.Args[1] {
	case "-enc":
		ckey := getKey()
		fi, _ := os.Open(os.Args[2]) 
		read := bufio.NewReader(fi)
		data, _ := ioutil.ReadAll(read) 
		fo, _ := os.Create(os.Args[2]+".enc") 
		fo.Write(encrypt(data, ckey))
		fmt.Println("done")
	case "-dec":
		ckey := getKey()
		fi, _ := os.Open(os.Args[2]) 
		read := bufio.NewReader(fi)
		data, _ := ioutil.ReadAll(read) 
		fo, _ := os.Create(os.Args[2]+".dec") 
		fo.Write(decrypt(data,ckey))
		fmt.Println("done")
	default:
		fmt.Println("usage: "+os.Args[0]+" [-enc|-dec] <file>")
	}	
}


