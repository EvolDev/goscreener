package storage

import "os"

func SaveScreenshot(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}
