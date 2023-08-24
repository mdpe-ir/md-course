package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	// Get image path from user
	fmt.Print("Enter the image path: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	imagePath := scanner.Text()

	// Read image file into byte slice
	data, err := ioutil.ReadFile(imagePath)
	if err != nil {
		log.Fatal(err)
	}

	// Get the base name of the input image without extension
	imageBaseName := strings.TrimSuffix(imagePath, ".jpg")

	// Loop through EXIF data and convert endianness
	for i := 0; i < len(data); i++ {
		if i >= 6 && string(data[i:i+5]) == "Exif\x00" {
			data = convertExifByteOrder(data, i+6)
			break
		}
	}

	// Generate the output image name based on the input image name
	outputImagePath := fmt.Sprintf("%s-modified.jpg", imageBaseName)

	// Write modified image back to file
	err = ioutil.WriteFile(outputImagePath, data, 0644)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Modified image saved as %s\n", outputImagePath)
}

func convertExifByteOrder(data []byte, offset int) []byte {
	// Read exif byte order marker
	order := binary.LittleEndian.Uint16(data[offset:])

	// If little endian, convert to big endian
	if order == 0x4949 {
		binary.BigEndian.PutUint16(data[offset:], 0x4D4D)
	}

	// If big endian, convert to little endian
	if order == 0x4D4D {
		binary.LittleEndian.PutUint16(data[offset:], 0x4949)
	}

	return data
}
