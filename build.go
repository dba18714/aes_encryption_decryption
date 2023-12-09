package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	buildAll := os.Getenv("BUILD_ALL_PLATFORMS") == "true"

	if buildAll {
		buildForAllPlatforms()
	} else {
		performInteractiveBuild()
	}
}

func buildForAllPlatforms() {
	configs := []struct {
		goos   string
		goarch string
		output string
	}{
		{"windows", "amd64", "aes_for_Windows64.exe"},
		{"windows", "386", "aes_for_Windows32.exe"},
		{"darwin", "amd64", "aes_for_macOS_AMD64"},
		{"darwin", "arm64", "aes_for_macOS_ARM64"},
		{"linux", "amd64", "aes_for_Linux64"},
		{"linux", "386", "aes_for_Linux32"},
		{"linux", "arm", "aes_for_Linux_ARM"},
	}

	for _, cfg := range configs {
		fmt.Printf("正在编译: %s\n", cfg.output)
		buildAndCompile(cfg.goos, cfg.goarch, filepath.Join("build", cfg.output))
	}
	fmt.Println("所有平台编译完成.")
}

func performInteractiveBuild() {
	fmt.Println("请选择目标平台:")
	fmt.Println("1. Windows 64位")
	fmt.Println("2. Windows 32位")
	fmt.Println("3. macOS AMD64")
	fmt.Println("4. macOS ARM64")
	fmt.Println("5. Linux 64位")
	fmt.Println("6. Linux 32位")
	fmt.Println("7. Linux ARM")
	fmt.Println("8. 编译以上所有平台")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("输入选择 (1-8): ")
	input, _ := reader.ReadString('\n')
	choice, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil {
		fmt.Println("无效输入")
		os.Exit(1)
	}

	var goos, goarch, output string
	buildDir := "build/"

	switch choice {
	case 1:
		goos = "windows"
		goarch = "amd64"
		output = "aes_for_Windows64.exe"
	case 2:
		goos = "windows"
		goarch = "386"
		output = "aes_for_Windows32.exe"
	case 3:
		goos = "darwin"
		goarch = "amd64"
		output = "aes_for_macOS_AMD64"
	case 4:
		goos = "darwin"
		goarch = "arm64"
		output = "aes_for_macOS_ARM64"
	case 5:
		goos = "linux"
		goarch = "amd64"
		output = "aes_for_Linux64"
	case 6:
		goos = "linux"
		goarch = "386"
		output = "aes_for_Linux32"
	case 7:
		goos = "linux"
		goarch = "arm"
		output = "aes_for_Linux_ARM"
	case 8:
		buildForAllPlatforms()
		return
	default:
		fmt.Println("无效选择")
		os.Exit(1)
	}

	buildAndCompile(goos, goarch, filepath.Join(buildDir, output))
	fmt.Printf("成功编译: %s\n", filepath.Join(buildDir, output))
}

func buildAndCompile(goos, goarch, output string) {
	cmd := exec.Command("go", "build", "-o", output, "aes_encryption_decryption.go")
	cmd.Env = append(os.Environ(), "GOOS="+goos, "GOARCH="+goarch)
	if err := cmd.Run(); err != nil {
		fmt.Printf("编译错误: %s\n", err)
		os.Exit(1)
	}
}
