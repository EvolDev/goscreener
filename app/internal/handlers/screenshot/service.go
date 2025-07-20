package screenshot

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"goscreener/internal/jsfunctions"
	"goscreener/internal/model"
	"goscreener/internal/retry"
	"io"
	"net/http"
	"strconv"
	"time"
)

func (h *Handler) TakeScreenshot(ctx context.Context, params model.ScreenshotParams, url string) ([]byte, error) {
	retryer := retry.NewRetryInterceptor(3)
	return retryer.DoWithRetry(ctx, url,
		func(context.Context) ([]byte, int, error) {
			screenshot, httpCode, err := captureScreenshot(ctx, params, url)
			return screenshot, httpCode, err
		})
}

func captureScreenshot(ctx context.Context, params model.ScreenshotParams, url string) ([]byte, int, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", false),
	)

	allocCtx, allocCancel := chromedp.NewExecAllocator(ctx, opts...)
	defer allocCancel()

	// Error log in console - off
	chromeCtx, chromeCancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(func(format string, a ...interface{}) {}))
	defer chromeCancel()

	timeoutCtx, cancel := context.WithTimeout(chromeCtx, 90*time.Second) // Или нужное тебе значение
	defer cancel()

	var httpStatusCode int
	var httpStatusText string
	var mainRequestID network.RequestID
	var screenshot []byte
	var elementBottom float64
	width, height, quality := params.GetDimensions()
	loadPageTimeout := params.LoadPageTimeoutSeconds
	finalLoadPageTimeout := params.FinalLoadPageTimeoutSeconds
	targetSelector := params.TargetSelector

	// We catch the main request to see the response code,
	// it is used later to send a repeat request in case of a non-200 response
	chromedp.ListenTarget(chromeCtx, func(ev interface{}) {
		if ev, ok := ev.(*network.EventRequestWillBeSent); ok {
			if ev.Type == "Document" && ev.Request.URL == url {
				mainRequestID = ev.RequestID
			}
		}
		if ev, ok := ev.(*network.EventResponseReceived); ok {
			if ev.RequestID == mainRequestID {
				httpStatusCode = int(ev.Response.Status)
				httpStatusText = ev.Response.StatusText
			}
		}
	})

	fakeNavHTMLStr := ""
	if params.FakeNav {
		fakeNavURL := "http://localhost:8080/fake-nav"
		resp, err := http.Get(fakeNavURL)
		if err != nil {
			return nil, http.StatusNotFound, fmt.Errorf("failed to fetch fake-nav HTML: %v", err)
		}
		defer resp.Body.Close()
		fakeNavHTML, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, http.StatusNotFound, fmt.Errorf("failed to read fake-nav HTML: %v", err)
		}
		// Convert fakeNavHTML to a string
		fakeNavHTMLStr = string(fakeNavHTML)
	}

	err := chromedp.Run(timeoutCtx, chromedp.Tasks{
		chromedp.EmulateViewport(width, height),
		chromedp.Navigate(url),
		chromedp.ActionFunc(func(ctx context.Context) error {
			if httpStatusCode == http.StatusNotFound {
				return fmt.Errorf("HTTP 404, page not found")
			}
			if httpStatusCode != http.StatusOK {
				return fmt.Errorf("HTTP %d, %s", httpStatusCode, httpStatusText)
			}
			return nil
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			if loadPageTimeout != 0 {
				return SleepContext(ctx, loadPageTimeout)
			}
			return nil
		}),
		// Inject the fake-nav HTML at the top of the page
		chromedp.ActionFunc(func(ctx context.Context) error {
			if !params.FakeNav {
				return nil
			}
			script := fmt.Sprintf(`
        (function() {
            var div = document.createElement('div');
            div.innerHTML = %s;
            document.body.insertBefore(div, document.body.firstChild);
        })();
   		 `, strconv.Quote(fakeNavHTMLStr))
			return chromedp.Evaluate(script, nil).Do(ctx)
		}),
		// Parse site title and fill placeholders in fake nav
		chromedp.ActionFunc(func(ctx context.Context) error {
			if !params.FakeNav {
				return nil
			}
			script := jsfunctions.AttachTitle()
			return chromedp.Evaluate(script, nil).Do(ctx)
		}),
		// Parse site icon and fill placeholders in fake nav
		chromedp.ActionFunc(func(ctx context.Context) error {
			if !params.FakeNav {
				return nil
			}
			script := jsfunctions.AttachIcon()
			return chromedp.Evaluate(script, nil).Do(ctx)
		}),
		// Fixed elements position
		chromedp.ActionFunc(func(ctx context.Context) error {
			if params.FixedNodes != nil {
				jsFunctionsStr := jsfunctions.MakeFixedNodes(params.FixedNodes)
				if jsFunctionsStr != "" {
					return chromedp.Evaluate(jsFunctionsStr, nil).Do(ctx)
				}
			}
			return nil
		}),
		// Scroll down, then up, for fullscreen format
		chromedp.ActionFunc(func(ctx context.Context) error {
			if params.FullScreen && params.WithScroll {
				jsFunctionsStr := jsfunctions.MakeFullScreenActions()
				if jsFunctionsStr != "" {
					return chromedp.Evaluate(jsFunctionsStr, nil).Do(ctx)
				}
			}
			return nil
		}),
		// Remove not wanted nodes, their wrote in request params
		chromedp.ActionFunc(func(ctx context.Context) error {
			if params.RemoveNodes != nil {
				jsFunctionsStr := jsfunctions.MakeRemoveNodesScript(params.RemoveNodes)
				if jsFunctionsStr != "" {
					return chromedp.Evaluate(jsFunctionsStr, nil).Do(ctx)
				}
			}
			return nil
		}),
		// Get bottom border from element
		chromedp.ActionFunc(func(ctx context.Context) error {
			if targetSelector != "" {
				script := fmt.Sprintf(`
                    (function() {
                        var elem = document.querySelector(%q);
                        if (!elem) {
                            return -1;
                        }
                        var rect = elem.getBoundingClientRect();
                        return rect.bottom + window.scrollY;
                    })();
                `, targetSelector)
				return chromedp.Evaluate(script, &elementBottom).Do(ctx)
			}
			return nil
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			if targetSelector != "" && elementBottom > 0 {
				// Set height for view range equal bottom element border
				return chromedp.EmulateViewport(width, int64(elementBottom)).Do(ctx)
			}
			return nil
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			if finalLoadPageTimeout != 0 {
				return SleepContext(ctx, finalLoadPageTimeout)
			}
			return nil
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			if loadPageTimeout != 0 {
				return SleepContext(ctx, loadPageTimeout)
			}
			return nil
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			if params.FullScreen {
				return chromedp.FullScreenshot(&screenshot, quality).Do(ctx)
			} else {
				return chromedp.CaptureScreenshot(&screenshot).Do(ctx)
			}
		}),
	})

	if err != nil {
		return nil, httpStatusCode, err
	}

	return screenshot, httpStatusCode, nil
}
