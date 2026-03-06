package coupon

import (
	"compress/gzip"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateCoupon_LengthBounds(t *testing.T) {
	svc := NewCouponService()
	assert.False(t, svc.ValidateCoupon("ABC"))
	assert.False(t, svc.ValidateCoupon("ABCDEFGHIJK"))
}

func TestValidateCoupon_FoundInTwoFiles(t *testing.T) {
	dir := t.TempDir()
	a := filepath.Join(dir, "a.gz")
	b := filepath.Join(dir, "b.gz")
	c := filepath.Join(dir, "c.gz")

	require.NoError(t, writeGz(a, []string{"HAPPYHRS"}))
	require.NoError(t, writeGz(b, []string{"HAPPYHRS"}))
	require.NoError(t, writeGz(c, []string{"FIFTYOFF"}))

	svc := NewCouponService(a, b, c)
	assert.True(t, svc.ValidateCoupon("HAPPYHRS"))
	assert.False(t, svc.ValidateCoupon("FIFTYOFF"))
}

func TestValidateCoupon_FileErrors(t *testing.T) {
	dir := t.TempDir()
	missing := filepath.Join(dir, "missing.gz")

	invalid := filepath.Join(dir, "invalid.gz")
	require.NoError(t, os.WriteFile(invalid, []byte("not gzipped"), 0o644))

	svc := NewCouponService(missing, invalid)
	assert.False(t, svc.ValidateCoupon("HAPPYHRS"))
}

func writeGz(path string, lines []string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	gz := gzip.NewWriter(f)
	for _, line := range lines {
		if _, err := gz.Write([]byte(line + "\n")); err != nil {
			_ = gz.Close()
			return err
		}
	}
	return gz.Close()
}

