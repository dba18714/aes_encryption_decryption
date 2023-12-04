package main

import (
	"bufio"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var choice string

	for {
		// 询问用户是要加密还是解密
		fmt.Printf("\n\u001B[32m请选择操作: \u001B[0m\n")
		fmt.Println("1: 加密")
		fmt.Println("2: 解密")
		fmt.Print("请输入选择 (1 或 2): ")
		choice, _ = reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		if choice == "1" || choice == "2" {
			break
		}

		fmt.Println("\u001B[31m无效的选择，请输入 1 或 2\u001B[0m")
	}

	// 输入密钥
	if choice == "1" {
		fmt.Printf("\n\u001B[32m请设置密钥（32个字符）: \u001B[0m\n")
	} else {
		fmt.Printf("\n\u001B[32m请输入密钥（32个字符）: \u001B[0m\n")
	}
	key, _ := reader.ReadString('\n')
	key = strings.TrimSpace(key)
	if len(key) != 32 {
		fmt.Println("\u001B[31m密钥长度必须是32个字符\u001B[0m")
		return
	}

	if choice == "1" {
		encryptData(reader, key)
	} else {
		decryptData(reader, key)
	}
}

func encryptData(reader *bufio.Reader, key string) {
	fmt.Printf("\n\u001B[32m请输入要加密的数据，输入'###'结束：\u001B[0m\n")

	var lines []string
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "读取错误:", err)
			break
		}
		line = strings.TrimSuffix(line, "\n")
		if line == "###" {
			break
		}
		lines = append(lines, line)
	}

	data := strings.Join(lines, "\n")

	// 检查数据是否为空
	if data == "" {
		fmt.Println("\u001B[31m错误：要加密的数据不能为空\u001B[0m")
		return
	}

	// 加密过程
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err)
	}

	// 生成 IV
	iv := make([]byte, block.BlockSize())
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	cbc := cipher.NewCBCEncrypter(block, iv)
	padding := block.BlockSize() - len(data)%block.BlockSize()
	paddedData := append([]byte(data), bytes.Repeat([]byte{byte(padding)}, padding)...)
	encrypted := make([]byte, len(paddedData))
	cbc.CryptBlocks(encrypted, paddedData)

	encryptedDataWithIv := append(iv, encrypted...)
	fmt.Printf("\n\u001B[32m加密后的数据（base64格式，包含IV）: \u001B[0m\n%v\n\n", base64.StdEncoding.EncodeToString(encryptedDataWithIv))
}

func decryptData(reader *bufio.Reader, key string) {
	fmt.Printf("\n\u001B[32m请输入密文（base64格式，包含IV）: \u001B[0m\n")
	encryptedBase64, _ := reader.ReadString('\n')
	encryptedBase64 = strings.TrimSpace(encryptedBase64)

	encryptedDataWithIv, err := base64.StdEncoding.DecodeString(encryptedBase64)
	if err != nil {
		fmt.Println("\u001B[32mbase64解码错误:\u001B[0m", err)
		return
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err)
	}

	if len(encryptedDataWithIv) < block.BlockSize() {
		fmt.Println("\u001B[32m加密数据太短\u001B[0m")
		return
	}

	iv := encryptedDataWithIv[:block.BlockSize()]
	encryptedData := encryptedDataWithIv[block.BlockSize():]

	cbcDecrypter := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(encryptedData))
	cbcDecrypter.CryptBlocks(decrypted, encryptedData)

	// 移除填充
	paddingSize := decrypted[len(decrypted)-1]
	decrypted = decrypted[:len(decrypted)-int(paddingSize)]

	fmt.Printf("\n\u001B[32m解密后的原始数据: \u001B[0m\n%v\n\n", string(decrypted))
}
