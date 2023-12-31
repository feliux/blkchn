package utils

import (
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strconv"
	"time"
)

var PATTERN = regexp.MustCompile(`((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?\.){3})(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)`)

func GetHost() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "127.0.0.1" // change depending of network environment
	}
	address, err := net.LookupHost(hostname)
	if err != nil {
		return "127.0.0.1" // change depending of network environment
	}
	return address[0]
}

// func IsFoundHost(host string, port uint16) bool {
func IsFoundHost(host string, port int) bool {
	target := fmt.Sprintf("%s:%d", host, port)
	_, err := net.DialTimeout("tcp", target, 1*time.Second)
	if err != nil {
		log.Printf("ERROR: %s", err)
		return false
	}
	log.Printf("TCP connection to %s. Added new host to neighbors.", target)
	return true
}

// func FindNeighbors(myHost string, myPort uint16, startIp uint8, endIp uint8, startPort uint16, endPort uint16) []string {
func FindNeighbors(myHost string, myPort int, startIp int, endIp int, startPort int, endPort int) []string {
	address := fmt.Sprintf("%s:%d", myHost, myPort)
	m := PATTERN.FindStringSubmatch(myHost)
	//log.Printf("Matched prefixHost %s", m) // [127.0.1.1 127.0.1. 1. 1]
	if m == nil {
		return nil
	}
	prefixHost := m[1]
	lastIp, _ := strconv.Atoi(m[len(m)-1])
	neighbors := make([]string, 0)
	for port := startPort; port <= endPort; port += 1 {
		for ip := startIp; ip <= endIp-1; ip += 1 {
			guessHost := fmt.Sprintf("%s%d", prefixHost, lastIp+ip)
			guessTarget := fmt.Sprintf("%s:%d", guessHost, port)
			log.Printf("Checking connection with host %s", guessTarget)
			if guessTarget != address && IsFoundHost(guessHost, port) {
				neighbors = append(neighbors, guessTarget)
			}
		}
	}
	log.Printf("Active blockchain neighbors: %v", neighbors)
	return neighbors
}
