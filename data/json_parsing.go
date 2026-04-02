package data

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type ThemesSongs struct {
	Themes []struct {
		ThemeName string   `json:"theme_name"`
		Songs     []string `json:"songs"`
	}
}

var (
	ErrInvalidInput      = errors.New("invalid input, song doesn't have any theme")
	ErrInvalidBytesInput = errors.New("Bytes cannot be null")
	CachedData           = make(map[string]string)
	Mu                   sync.RWMutex
)

func LoadData(filename string) error {

	fileBytes, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	if err := swapData(fileBytes); err != nil {
		return err
	}

	return nil

}

func swapData(data []byte) error {
	var loadedData ThemesSongs

	if err := json.Unmarshal(data, &loadedData); err != nil {
		return err
	}

	var swapNewData = make(map[string]string)
	for _, themeMedata := range loadedData.Themes {
		for _, song := range themeMedata.Songs {
			swapNewData[song] = themeMedata.ThemeName
		}
	}
	Mu.Lock()
	CachedData = swapNewData
	Mu.Unlock()

	return nil
}

func FetchTheme(name string) (string, error) {

	var err error

	Mu.RLock()
	value, ok := CachedData[name]
	Mu.RUnlock()

	if !ok {
		err = ErrInvalidInput
	}

	return value, err
}
