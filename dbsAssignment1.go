package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func processHeader(headerText string) map[string]string {
	// Initialize Variables
	binCharMap := make(map[string]string)
	binMarker := "1"
	binCounter := 0

	for _, header := range headerText {
		// Format binary counter to base 2
		binText := strconv.FormatInt(int64(binCounter), 2)

		// Check if converted counter matches withj
		// binary marker (contains only "1").
		// If binary text matches binary marker,
		// append binary counter, to replicate skipping counter
		// and convert next counter
		if binText == binMarker {
			binMarker += "1"
			binText = strings.Replace(binText, "1", "0", -1) + "0"
			binCounter = 0
		}

		// Maps converted counter to header character
		// Pads "0" left of string according to length of marker
		binCharMap[fmt.Sprintf("%0"+strconv.Itoa(len(binMarker))+"s", binText)] = string(header)
		binCounter++
	}
	return binCharMap
}

func processMessage(messageText string, binCharMap map[string]string) string {
	decodedMessage := ""
	messageProcessing := true
	for messageProcessing {
		// Get length of segment (First 3 binary)
		// Convert to integer
		segLen, _ := strconv.ParseInt(string(messageText[0:3]), 2, 64)

		// If segment length is more than 0
		// start processing subsequent segment
		if segLen > 0 {
			// Remove converted binary which was used
			// to find the length of segment
			messageText = messageText[3:]

			// Loops segment until a segment with all "1" is found
			segProcessing := true
			for segProcessing {
				// Declare regular expression to find "0"
				segMarker := regexp.MustCompile("[0]+")

				// If no "0" is found in segment, end loop
				if !segMarker.MatchString(messageText[0:segLen]) {
					segProcessing = false
				} else {
					// Maps segment binary to binary key in binary character map
					// and print out the value.
					decodedMessage += binCharMap[messageText[0:segLen]]
				}
				// Remove segment from body text
				// once processing is done on segment
				messageText = messageText[segLen:]
			}
		} else {
			messageProcessing = false
		}
	}
	return decodedMessage
}

func main() {
	args := os.Args
	// Opens file for processing
	file, err := os.Open(args[1])
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	// Initialize variables
	// Initialize Scanner and Read File
	// Default uses new line as token
	scanner := bufio.NewScanner(file)

	headerReceived, messageReceived := false, false
	headerText, messageText := "", ""
	nextScan := true

	for nextScan {
		nextScan = scanner.Scan()

		// Checks line length,
		// skips line if no length
		if len(scanner.Text()) > 0 {
			if !headerReceived {
				// Receive header text
				headerText = scanner.Text()
				headerReceived = true
			} else {
				// Receive message text, stops receiving when last 3 characters are "000"
				messageText += strings.Replace(scanner.Text(), " \r\n", "", -1)

				if string(messageText[len(messageText)-3:]) == "000" {
					messageReceived = true
				}
			}
		}

		// Start processing header and message
		// Output decoded message
		if headerReceived && messageReceived {
			binCharMap := processHeader(headerText)
			decodedMessage := processMessage(messageText, binCharMap)
			headerReceived, messageReceived = false, false
			headerText, messageText = "", ""
			fmt.Println(decodedMessage)
		}
	}
}
