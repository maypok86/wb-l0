package cache

import (
	"testing"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/maypok86/wb-l0/internal/entity"
	"github.com/stretchr/testify/require"
)

func TestNewMemoryCache(t *testing.T) {
	c := NewMemoryCache()
	require.NotNil(t, c)
}

func fakeOrder(t *testing.T) *entity.Order {
	t.Helper()

	order := &entity.Order{}
	require.NoError(t, faker.FakeData(order))
	return order
}

func TestMemoryCache_Get(t *testing.T) {
	c := NewMemoryCache()
	order := fakeOrder(t)
	gotOrder, err := c.Get(0)
	require.Error(t, err)
	require.Nil(t, gotOrder)
	c.cache[0] = &item{
		value:     order,
		createdAt: time.Now().Unix(),
		ttl:       61,
	}
	gotOrder, err = c.Get(0)
	require.NoError(t, err)
	require.Equal(t, order, gotOrder)
}

func TestMemoryCache_Set(t *testing.T) {
	c := NewMemoryCache()
	order := fakeOrder(t)
	require.NoError(t, c.Set(0, order, 1))
	gotOrder, err := c.Get(0)
	require.NoError(t, err)
	require.Equal(t, order, gotOrder)
}

func TestMemoryCache_Clean(t *testing.T) {
	c := NewMemoryCache()
	order := fakeOrder(t)
	require.NoError(t, c.Set(0, order, 1))
	gotOrder, err := c.Get(0)
	require.NoError(t, err)
	require.Equal(t, order, gotOrder)
	time.Sleep(2 * time.Second)
	gotOrder, err = c.Get(0)
	require.Error(t, err)
	require.Nil(t, gotOrder)
}
