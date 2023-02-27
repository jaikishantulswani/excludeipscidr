package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
)

func main() {
	cidrFileFlag := flag.String("c", "", "File containing CIDR ranges to exclude")
	fileFlag := flag.String("file", "", "File to read IP addresses from")
	flag.Parse()

	var cidrs []*net.IPNet
	if *cidrFileFlag != "" {
		cidrFile, err := os.Open(*cidrFileFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening CIDR file: %v\n", err)
			os.Exit(1)
		}
		defer cidrFile.Close()

		scanner := bufio.NewScanner(cidrFile)
		for scanner.Scan() {
			cidrString := scanner.Text()

			_, cidr, err := net.ParseCIDR(cidrString)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error parsing CIDR range: %v\n", err)
				continue
			}

			cidrs = append(cidrs, cidr)
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading CIDR file: %v\n", err)
			os.Exit(1)
		}
	}

	var file *os.File
	if *fileFlag != "" {
		var err error
		file, err = os.Open(*fileFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()
	} else {
		file = os.Stdin
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ipString := scanner.Text()

		ip := net.ParseIP(ipString)
		if ip == nil {
			fmt.Fprintf(os.Stderr, "Invalid IP address: %v\n", ipString)
			continue
		}

		excluded := false
		for _, cidr := range cidrs {
			if cidr.Contains(ip) {
				excluded = true
				break
			}
		}

		if excluded {
			continue
		}

		fmt.Println(ipString)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}
}
