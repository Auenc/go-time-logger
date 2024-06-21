package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func loadFileContents(path string) ([]byte, error) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	var bs []byte
	for {
		b, err := reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		bs = append(bs, b)
	}
	return bs, nil
}

func writeFileContents(path string, data []byte) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("error opening file %w", err)
	}
	defer f.Close()

	err = f.Truncate(0)
	if err != nil {
		return fmt.Errorf("error truncating file: %w", err)
	}

	_, err = f.Write(data)
	if err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}

	return nil
}
