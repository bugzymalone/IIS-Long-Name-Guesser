package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

// Color codes for console output
const (
	Reset = "\033[0m"
	Red   = "\033[31m"
	Green = "\033[32m"
	Blue  = "\033[34m"
)

func main() {
	urlPtr := flag.String("url", "", "URL to test")
	wordStubPtr := flag.String("stub", "", "Word stub (optional)")
	wordExtrapolatePtr := flag.String("extrapolate", "", "Word to extrapolate")
	fileExtensionPtr := flag.String("extension", "", "File extension")
	dictionaryPathPtr := flag.String("dictionary", "", "Dictionary to use")
	suffixPathPtr := flag.String("suffix", "", "Suffix file")

	flag.Parse()

	if *urlPtr == "" || *wordExtrapolatePtr == "" || *fileExtensionPtr == "" || *dictionaryPathPtr == "" {
		log.Fatal("IIS Long Name Guesser. Command Line Arguments:\n\n--url\n\n--stub (optional)\n\n--extrapolate\n\n--suffix (optional)\n\n--extension\n\n--dictionary\n\nE.g., to guess the long version of 'bigque~1.aspx' where the second word component starts with 'que' from the dictionary, use:\n\n--url=https://www.example.com/ --stub=big --extrapolate=que --suffix=suffix.txt --extension=aspx --dictionary=dict.txt")
	}

	dictionaryFile, err := os.Open(*dictionaryPathPtr)
	if err != nil {
		log.Fatalf("Failed to open dictionary file: %s", err)
	}
	defer dictionaryFile.Close()

	var suffixes []string
	if *suffixPathPtr != "" {
		suffixFile, err := os.Open(*suffixPathPtr)
		if err != nil && !os.IsNotExist(err) {
			log.Fatalf("Failed to open suffix file: %s", err)
		} else if err == nil {
			defer suffixFile.Close()

			suffixScanner := bufio.NewScanner(suffixFile)
			for suffixScanner.Scan() {
				suffix := suffixScanner.Text()
				suffixes = append(suffixes, suffix)
			}
			if err := suffixScanner.Err(); err != nil {
				log.Fatalf("Failed to read suffix file: %s", err)
			}
		}
	}

	var words []string
	scanner := bufio.NewScanner(dictionaryFile)
	for scanner.Scan() {
		word := scanner.Text()
		if strings.HasPrefix(word, *wordExtrapolatePtr) {
			words = append(words, word)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Failed to read dictionary file: %s", err)
	}

	for _, word := range words {
		if len(suffixes) > 0 {
			for _, suffix := range suffixes {
				finalURL := fmt.Sprintf("%s%s%s%s.%s", *urlPtr, *wordStubPtr, word, suffix, *fileExtensionPtr)
				response, err := http.Get(finalURL)
				if err != nil {
					log.Printf("HTTP request failed: %s", err)
					continue
				}

				statusColor := Reset // Default color
				switch response.StatusCode {
				case 200:
					statusColor = Green
					log.Printf("URL: %s, Response code: %s%d%s, Body size: %d\n", finalURL, statusColor, response.StatusCode, Reset, response.ContentLength)
					response.Body.Close()
					return // Exit the program when 200 response is encountered
				case 404:
					statusColor = Blue
				case 403:
					statusColor = Red
				}

				log.Printf("URL: %s, Response code: %s%d%s, Body size: %d\n", finalURL, statusColor, response.StatusCode, Reset, response.ContentLength)
				response.Body.Close()
			}
		} else {
			finalURL := fmt.Sprintf("%s%s%s.%s", *urlPtr, *wordStubPtr, word, *fileExtensionPtr)
			response, err := http.Get(finalURL)
			if err != nil {
				log.Printf("HTTP request failed: %s", err)
				continue
			}

			statusColor := Reset // Default color
			switch response.StatusCode {
			case 200:
				statusColor = Green
				log.Printf("URL: %s, Response code: %s%d%s, Body size: %d\n", finalURL, statusColor, response.StatusCode, Reset, response.ContentLength)
				response.Body.Close()
				return // Exit the program when 200 response is encountered
			case 404:
				statusColor = Blue
			case 403:
				statusColor = Red
			}

			log.Printf("URL: %s, Response code: %s%d%s, Body size: %d\n", finalURL, statusColor, response.StatusCode, Reset, response.ContentLength)
			response.Body.Close()
		}
	}
}
