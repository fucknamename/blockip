package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	/*
		IP 列表文件 ip_list.txt 示例：
			192.168.1.1
			192.168.2.0/24
			10.0.0.1
	*/
	ipFile := "ip.txt"    // IP 列表文件
	ruleName := "blockip" // 已创建的规则名称

	// 读取 IP 列表
	ips, err := readIPList(ipFile)
	if err != nil {
		fmt.Printf("Failed to read IP list: %v\n", err)
		return
	}

	// // 批量添加防火墙规则
	// for _, ip := range ips {
	// 	err := blockIP(ip)
	// 	if err != nil {
	// 		log.Printf("Failed to block IP %s: %v", ip, err)
	// 	} else {
	// 		fmt.Printf("Blocked IP: %s\n", ip)
	// 	}
	// }

	// 移除指定IP，可以先把规则移除，然后再添加规则，批量添加IP进去

	// 更新防火墙规则
	err = updateFirewallRule(ruleName, ips)
	if err != nil {
		fmt.Printf("Failed to update firewall rule: %v\n", err)
	} else {
		fmt.Println("Firewall rule updated successfully!")
	}
}

func readIPList(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var ips []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ips = append(ips, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return ips, nil
}

// 更新防火墙规则，批量添加 IP
func updateFirewallRule(ruleName string, ips []string) error {
	// 将 IP 列表拼接为逗号分隔的字符串
	ipList := strings.Join(ips, ",")
	cmd := exec.Command("netsh", "advfirewall", "firewall", "set", "rule",
		"name="+ruleName, "new", "remoteip="+ipList)
	return cmd.Run()
}

// 添加防火墙规则
func blockIP(ruleName, ip string) error {
	cmd := exec.Command("netsh", "advfirewall", "firewall", "add", "rule",
		"name="+ruleName, "dir=in", "action=block", "remoteip="+ip, "enable=yes")
	return cmd.Run()
}

func removeBlockRule(ruleName string) error {
	cmd := exec.Command("netsh", "advfirewall", "firewall", "delete", "rule",
		"name="+ruleName)
	return cmd.Run()
}

func isAdmin() bool {
	return os.Getenv("USERNAME") == "Administrator"
}
