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

func main() {
	urlPtr := flag.String("url", "", "URL to test")
	wordStubPtr := flag.String("stub", "", "Word stub (optional)")
	wordExtrapolatePtr := flag.String("extrapolate", "", "Word to extrapolate")
	fileExtensionPtr := flag.String("extension", "", "File extension")
	dictionaryPathPtr := flag.String("dictionary", "", "Dictionary to use")

	flag.Parse()

	if *urlPtr == "" || *wordExtrapolatePtr == "" || *fileExtensionPtr == "" || *dictionaryPathPtr == "" {
		log.Fatal("\n IIS Long Name Guesser \n\n E.G --url=http:\\\\127.0.0.1/ --stub=big (optional arg) --extrapolate=vis --extension=aspx --dictionary=words.txt")
	}

	file, err := os.Open(*dictionaryPathPtr)
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var words []string
	for scanner.Scan() {
		word := scanner.Text()
		if strings.HasPrefix(word, *wordExtrapolatePtr) {
			words = append(words, word)
		}
	}

	file.Close()

	for _, word := range words {
		finalUrl := fmt.Sprintf("%s%s%s.%s", *urlPtr, *wordStubPtr, word, *fileExtensionPtr)
		response, err := http.Get(finalUrl)
		if err != nil {
			log.Printf("HTTP request failed: %s", err)
		} else {
			log.Printf("URL: %s, Response code: %d, Body size: %d\n", finalUrl, response.StatusCode, response.ContentLength)
		}
	}
}
