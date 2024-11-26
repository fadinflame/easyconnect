package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type Config struct {
	Server    string `json:"server"`
	Group     string `json:"group"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	CiscoLogs bool   `json:"cisco_logs"`
}

func getConfigPath() (string, error) {
	var configDir string
	if runtime.GOOS == "windows" {
		configDir = filepath.Join(os.Getenv("APPDATA"), "easyconnect")
	} else {
		configDir = filepath.Join(os.Getenv("HOME"), ".easyconnect")
	}

	if err := os.MkdirAll(configDir, 0700); err != nil {
		return "", err
	}

	return filepath.Join(configDir, "config.json"), nil
}

func generateConfig() (*Config, error) {
	fmt.Println("Configuration file not found. Let's create one.")
	var config Config

	fmt.Print("Enter VPN server address: ")
	if _, err := fmt.Scanln(&config.Server); err != nil {
		return nil, err
	}

	fmt.Print("Enter VPN group name(number): ")
	if _, err := fmt.Scanln(&config.Group); err != nil {
		return nil, err
	}

	fmt.Print("Enter VPN username: ")
	if _, err := fmt.Scanln(&config.Username); err != nil {
		return nil, err
	}

	fmt.Print("Enter VPN password: ")
	if _, err := fmt.Scanln(&config.Password); err != nil {
		return nil, err
	}

	fmt.Print("Do you want to enable Cisco logs? (y/n): ")
	var ciscoLogs string
	if _, err := fmt.Scanln(&ciscoLogs); err != nil {
		return nil, err
	}
	ciscoLogs = strings.ToLower(ciscoLogs)

	if ciscoLogs == "y" {
		config.CiscoLogs = true
	}

	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	file, err := os.Create(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(&config); err != nil {
		return nil, err
	}

	fmt.Printf("Configuration saved to %s\n", configPath)
	return &config, nil
}

func loadConfig() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	file, err := os.Open(configPath)
	if os.IsNotExist(err) {
		return generateConfig()
	} else if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &Config{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(config); err != nil {
		return nil, err
	}

	if config.Server == "" || config.Group == "" || config.Username == "" || config.Password == "" {
		return generateConfig()
	}

	return config, nil
}

func executeCommand(command string, args []string, input string) (string, error) {
	cmd := exec.Command(command, args...)
	cmd.Stdin = bytes.NewBufferString(input)

	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output

	err := cmd.Run()
	return output.String(), err
}

func checkVpnStatus() (bool, error) {
	command := "/opt/cisco/anyconnect/bin/vpn"
	if runtime.GOOS == "windows" {
		command = "vpncli.exe"
	}

	output, err := executeCommand(command, []string{"status"}, "")
	if err != nil {
		return false, err
	}

	for _, line := range bytes.Split([]byte(output), []byte("\n")) {
		line = bytes.TrimSpace(line)
		if bytes.HasPrefix(line, []byte(">> state:")) {
			state := string(bytes.TrimPrefix(line, []byte(">> state:")))
			state = strings.TrimSpace(state)
			if strings.EqualFold(state, "Connected") {
				return true, nil
			}
		}
	}

	return false, nil
}

func connectToVpn(config *Config) error {
	command := "/opt/cisco/anyconnect/bin/vpn"
	if runtime.GOOS == "windows" {
		command = "vpncli.exe"
	}

	input := fmt.Sprintf("%s\n%s\n%s\ny\n", config.Group, config.Username, config.Password)
	output, err := executeCommand(command, []string{"-s", "connect", config.Server}, input)

	if config.CiscoLogs {
		fmt.Println(output)
	}

	return err
}

func disconnectFromVpn(config *Config) error {
	command := "/opt/cisco/anyconnect/bin/vpn"
	if runtime.GOOS == "windows" {
		command = "vpncli.exe"
	}

	output, err := executeCommand(command, []string{"disconnect"}, "")

	if config.CiscoLogs {
		fmt.Println(output)
	}

	return err
}

func containsIgnoreCase(source, substr string) bool {
	return bytes.Contains(bytes.ToLower([]byte(source)), bytes.ToLower([]byte(substr)))
}

func main() {
	config, err := loadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	isConnected, err := checkVpnStatus()
	if err != nil {
		fmt.Printf("Error checking VPN status: %v\n", err)
		os.Exit(1)
	}

	if isConnected {
		fmt.Println("VPN is currently connected. Disconnecting...")
		if err := disconnectFromVpn(config); err != nil {
			fmt.Printf("Error disconnecting VPN: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("VPN disconnected successfully.")
	} else {
		fmt.Println("VPN is not connected. Connecting...")
		if err := connectToVpn(config); err != nil {
			fmt.Printf("Error connecting VPN: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("VPN connected successfully.")
	}

	fmt.Println("Press Enter to exit...")
	fmt.Scanln()
}
