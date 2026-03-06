package coupon

import (
	"bufio"
	"compress/gzip"
	"os"
)

type CouponService struct {
	files []string
}

func NewCouponService(files ...string) *CouponService {
	return &CouponService{files: files}
}

func (c *CouponService) ValidateCoupon(code string) bool {

	// Coupon codes in this service are expected to be short, human-entered strings.
	// Project requirements include 7-char codes (e.g. "HAPPYHRS"), so keep the lower bound inclusive.
	if len(code) < 7 || len(code) > 10 {
		return false
	}

	results := make(chan bool)
	done := make(chan struct{})

	for _, file := range c.files {

		go func(f string) {

			select {
			case <-done:
				return
			default:
			}

			found := searchInFile(f, code)

			results <- found

		}(file)
	}

	count := 0

	for i := 0; i < len(c.files); i++ {

		if <-results {
			count++
		}

		if count >= 2 {
			close(done)
			return true
		}
	}

	return false
}

func searchInFile(path string, code string) bool {

	f, err := os.Open(path)
	if err != nil {
		return false
	}

	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return false
	}

	defer gz.Close()

	scanner := bufio.NewScanner(gz)

	for scanner.Scan() {

		if scanner.Text() == code {
			return true
		}
	}

	return false
}
