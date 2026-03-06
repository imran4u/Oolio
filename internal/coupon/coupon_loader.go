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

	if len(code) < 8 || len(code) > 10 {
		return false
	}

	count := 0

	for _, file := range c.files {

		found := searchInFile(file, code)

		if found {
			count++
		}

		if count >= 2 {
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
