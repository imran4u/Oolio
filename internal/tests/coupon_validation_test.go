package tests

import (
	"compress/gzip"
	"os"
	"path/filepath"
	"testing"

	"github.com/imran4u/Oolio/internal/coupon"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCouponService_ValidateCoupon(t *testing.T) {
	dir := t.TempDir()

	file1 := filepath.Join(dir, "a.gz")
	file2 := filepath.Join(dir, "b.gz")
	file3 := filepath.Join(dir, "c.gz")

	require.NoError(t, writeGzipLines(file1, []string{"HAPPYHRS", "NOPE0000"}))
	require.NoError(t, writeGzipLines(file2, []string{"HAPPYHRS"}))
	require.NoError(t, writeGzipLines(file3, []string{"FIFTYOFF"}))

	svc := coupon.NewCouponService(file1, file2, file3)

	assert.True(t, svc.ValidateCoupon("HAPPYHRS"), "present in 2 files")
	assert.False(t, svc.ValidateCoupon("FIFTYOFF"), "present in only 1 file")
	assert.False(t, svc.ValidateCoupon("SUPER100"), "not present in any file")

	assert.False(t, svc.ValidateCoupon("ABC"), "too short")
	assert.False(t, svc.ValidateCoupon("ABCDEFGHIJK"), "too long")
}

func TestCouponService_ValidateCoupon_ErrorPaths(t *testing.T) {
	dir := t.TempDir()

	// Non-existent file -> os.Open error path.
	missingPath := filepath.Join(dir, "missing.gz")

	// Invalid gzip -> gzip.NewReader error path.
	invalidGzip := filepath.Join(dir, "invalid.gz")
	require.NoError(t, os.WriteFile(invalidGzip, []byte("not gzipped"), 0o644))

	svc := coupon.NewCouponService(missingPath, invalidGzip)

	assert.False(t, svc.ValidateCoupon("HAPPYHRS"))
}

func writeGzipLines(path string, lines []string) error {
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

