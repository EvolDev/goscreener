package model

type ScreenshotParams struct {
	URL                         string          `json:"url"`
	URLs                        []string        `json:"urls"`
	LoadPageTimeoutSeconds      int             `json:"load_page_timeout_seconds"`
	FinalLoadPageTimeoutSeconds int             `json:"final_load_page_timeout_seconds"`
	FixedNodes                  []*NodeSelector `json:"fixed_nodes"`
	RemoveNodes                 []*NodeSelector `json:"remove_nodes"`
	Height                      int64           `json:"height"`
	Width                       int64           `json:"width"`
	Quality                     int             `json:"quality"`
	FullScreen                  bool            `json:"full_screen"`
	FakeNav                     bool            `json:"fake_nav"`
	WithScroll                  bool            `json:"with_scroll"`
	TargetSelector              string          `json:"target_selector"`
	Cache                       bool            `json:"cache"`
}

func (p *ScreenshotParams) GetDimensions() (int64, int64, int) {
	if p.Width == 0 {
		p.Width = 1080
	}
	if p.Height == 0 {
		p.Height = 1920
	}
	if p.Quality == 0 {
		p.Quality = 100
	}
	return int64(p.Width), int64(p.Height), p.Quality
}
