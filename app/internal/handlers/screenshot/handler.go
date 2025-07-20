package screenshot

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"goscreener/internal/model"
	"goscreener/internal/storage"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

const maxParallelRequests = 10

type Handler struct {
	host       string
	port       string
	screensDir string
}

func NewScreenshotHandler() *Handler {
	return &Handler{}
}

func NewScreenshotsManyHandler(host string, port string, screensDir string) http.HandlerFunc {
	return (&Handler{host: host, port: port, screensDir: screensDir}).ServeMultiple
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get URL from Request
	var screenshotParams model.ScreenshotParams
	err := json.NewDecoder(r.Body).Decode(&screenshotParams)
	//	w.Header().Set("Connection", "close")
	if err != nil {
		http.Error(w, "Params are required", http.StatusBadRequest)
		return
	}
	_ = r.Body.Close()

	if screenshotParams.URL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	// Screenshot capture
	screenshot, err := h.TakeScreenshot(ctx, screenshotParams, screenshotParams.URL)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to take screenshot: %v", err), http.StatusInternalServerError)
		return
	}

	// Set header Content-Type for image
	w.Header().Set("Content-Type", "image/png")
	// Write screenshot in response
	_, err = w.Write(screenshot)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to send screenshot: %v", err), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) ServeMultiple(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var screenshotParams model.ScreenshotParams
	err := json.NewDecoder(r.Body).Decode(&screenshotParams)
	//	w.Header().Set("Connection", "close")
	if err != nil {
		http.Error(w, "Params are required", http.StatusBadRequest)
		return
	}
	_ = r.Body.Close()

	if len(screenshotParams.URLs) == 0 {
		http.Error(w, "URLs are required", http.StatusBadRequest)
		return
	}

	var wg sync.WaitGroup
	mux := sync.Mutex{}
	semaphore := make(chan struct{}, maxParallelRequests)
	results := make(map[string]interface{})
	errors := make(map[string]string)

	for _, url := range screenshotParams.URLs {
		wg.Add(1)
		semaphore <- struct{}{}
		go func(url string) {
			defer func() {
				<-semaphore // clear semaphore
				wg.Done()
			}()

			select {
			case <-ctx.Done():
				fmt.Printf("Context canceled for %s\n", url)
				return
			default:
				screenshotParams.URL = url
				hash := md5.Sum([]byte(url))
				hashStr := hex.EncodeToString(hash[:])
				if screenshotParams.FullScreen {
					hashStr = fmt.Sprintf("%s%s", hashStr, "full")
				}
				filePath := filepath.Join(h.screensDir, fmt.Sprintf("%s.jpg", hashStr))
				fileURL := fmt.Sprintf("%s:%s/get/screenshot/%s.jpg", h.host, h.port, hashStr)

				// Search cached image
				if screenshotParams.Cache {
					if _, err := os.Stat(filePath); err == nil {
						// If it found, then return url
						mux.Lock()
						results[url] = fileURL
						mux.Unlock()
						return
					}
				}

				screenshot, err := h.TakeScreenshot(ctx, screenshotParams, url)
				if err != nil {
					mux.Lock()
					errors[url] = fmt.Sprintf("Error capturing screenshot for %s: %v", url, err)
					mux.Unlock()
					return
				}

				if err := storage.SaveScreenshot(filePath, screenshot); err != nil {
					fmt.Printf("Error saving screenshot for %s: %v\n", url, err)
					return
				}

				mux.Lock()
				results[url] = fileURL
				mux.Unlock()
			}
		}(url)
	}

	go func() {
		wg.Wait()
		close(semaphore)
	}()

	wg.Wait()

	w.Header().Set("Content-Type", "application/json")
	if len(errors) > 0 {
		mux.Lock()
		results["errors"] = errors
		mux.Unlock()
	}
	err = json.NewEncoder(w).Encode(results)
	if err != nil {
		return
	}
}
