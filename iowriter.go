package main

import (
	"bufio"
	"fmt"
	"os"
)

const maxLinesPerFile = 100000

func WriteFromChannel(ch <-chan string) {
	writer, err := NewFileWriter()
	if err != nil {
		fmt.Println(err)
		return
	}
	go func() {
		defer func(writer *FileWriter) {
			err := writer.Close()
			fmt.Println("Closing file writer")
			if err != nil {
				fmt.Println("Error closing file writer:", err)
			}
		}(writer)
		for line := range ch {
			err := writer.WriteLine(line)
			if err != nil {
				fmt.Println("Error writing line:", err)
				return
			}
		}
	}()
}

type FileWriter struct {
	file      *os.File
	writer    *bufio.Writer
	lineCount int
	fileIndex int
}

func NewFileWriter() (*FileWriter, error) {
	fw := &FileWriter{}

	// Find the latest file with space or create a new one
	for i := 1; ; i++ {
		fileName := fmt.Sprintf("translate%07d.txt", i)
		if _, err := os.Stat(fileName); os.IsNotExist(err) {
			// No more existing files, create a new one with the last found index
			fw.fileIndex = i - 1
			break
		} else {
			// File exists, check if it's full
			file, err := os.OpenFile(fileName, os.O_RDWR|os.O_APPEND, 0644)
			if err != nil {
				return nil, err
			}

			scanner := bufio.NewScanner(file)
			lineCount := 0
			for scanner.Scan() {
				lineCount++
			}

			file.Close() // Close the file after counting lines

			if lineCount < maxLinesPerFile {
				// File has space, use it
				fw.fileIndex = i
				fw.lineCount = lineCount
				fw.file, err = os.OpenFile(fileName, os.O_RDWR|os.O_APPEND, 0644)
				if err != nil {
					return nil, err
				}
				fw.writer = bufio.NewWriter(fw.file)
				return fw, nil
			}
		}
	}

	// Create the new file since no suitable existing file was found
	err := fw.rotateFile()
	if err != nil {
		return nil, err
	}
	return fw, nil
}

func (fw *FileWriter) rotateFile() error {
	if fw.writer != nil {
		fw.writer.Flush()
		fw.file.Close()
	}

	fw.fileIndex++
	fileName := fmt.Sprintf("translate%07d.txt", fw.fileIndex)
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}

	fw.file = file
	fw.writer = bufio.NewWriter(file)
	fw.lineCount = 0

	return nil
}

func (fw *FileWriter) WriteLine(line string) error {
	if fw.lineCount >= maxLinesPerFile {
		err := fw.rotateFile()
		if err != nil {
			return err
		}
	}

	_, err := fw.writer.WriteString(line + "\n")
	if err != nil {
		return err
	}

	fw.lineCount++
	return nil
}

func (fw *FileWriter) Close() error {
	if fw.writer != nil {
		err := fw.writer.Flush()
		if err != nil {
			return err
		}
	}

	if fw.file != nil {
		err := fw.file.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
