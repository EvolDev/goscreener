package retry

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"time"
)

type RetryInterceptor struct {
	maxRetries int
}

func NewRetryInterceptor(maxRetries int) *RetryInterceptor {
	return &RetryInterceptor{maxRetries: maxRetries}
}

func (ri *RetryInterceptor) DoWithRetry(ctx context.Context, url string, operation func(context.Context) ([]byte, int, error)) ([]byte, error) {
	//fmt.Println("Entering DoWithRetry")
	var result []byte
	var err error
	var httpStatusCode int

	for attempt := 0; attempt < ri.maxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(randomDelay(1, 3))
		}

		attemptCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
		defer cancel()

		now := time.Now()
		logTimeFormat := "2006-01-02 15:04:05"

		fmt.Printf("[%s] URL: %s, Attempt %d/%d\n", now.Format(logTimeFormat), url, attempt+1, ri.maxRetries)
		result, httpStatusCode, err = operation(attemptCtx)

		if err == nil && httpStatusCode == http.StatusOK {
			fmt.Printf("[%s] URL: %s, Operation succeeded\n", now.Format(logTimeFormat), url)
			return result, nil
		}

		fmt.Println("Operation failed with error:", err)

		//  Если ошибка 404 — выходим сразу
		if httpStatusCode == http.StatusNotFound {
			fmt.Printf("[%s] URL: %s, HTTP 404 received. Stopping retries.\n", now.Format(logTimeFormat), url)
			break
		}

		//  Прерываем, если ошибка не подлежит ретраю
		/*		if isNonRetryableError(err) {
				fmt.Printf("Non-retryable error encountered. Stopping retries. URL: %s\n", params.URL)
				break
			}*/
	}

	fmt.Println("All retry attempts exhausted.")
	return result, err
}

/*func isNonRetryableError(err error) bool {
	errMsg := err.Error()
	for _, nonRetryableError := range nonRetryableErrors {
		if strings.Contains(errMsg, nonRetryableError) {
			return true
		}
	}
	return false
}*/

func randomDelay(min, max int) time.Duration {
	delta := max - min
	randNum, _ := rand.Int(rand.Reader, big.NewInt(int64(delta+1)))
	return time.Duration(min+int(randNum.Int64())) * time.Second
}
