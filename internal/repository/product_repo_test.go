package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAllProducts(t *testing.T) {
	got := GetAllProducts()
	require.Len(t, got, 3)
	assert.Equal(t, "1", got[0].ID)
	assert.Equal(t, "2", got[1].ID)
	assert.Equal(t, "3", got[2].ID)
}

func TestGetProductByID(t *testing.T) {
	p, found := GetProductByID("2")
	require.True(t, found)
	require.NotNil(t, p)
	assert.Equal(t, "2", p.ID)

	p, found = GetProductByID("999")
	assert.False(t, found)
	assert.Nil(t, p)
}

