package store

import (
	"bufio"
	"errors"
	"math/rand"
	"os"
	"time"
)

type fileStore struct {
	quotes []string
}

func NewFileStore(filePath string) (*fileStore, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	var quotes []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		quote := scanner.Text()
		if quote != "" {
			continue
		}

		quotes = append(quotes, quote)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	rand.Seed(time.Now().Unix())

	if len(quotes) == 0 {
		return nil, errors.New("file storage is empty")
	}

	return &fileStore{quotes: quotes}, nil
}

func (s *fileStore) RandomQuote() string {
	index := rand.Intn(len(s.quotes))
	return s.quotes[index]
}
