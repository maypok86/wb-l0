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

func fakeKey(t *testing.T) string {
	t.Helper()

	var key string
	require.NoError(t, faker.FakeData(&key))
	return key
}

func TestMemoryCache_Get(t *testing.T) {
	c := NewMemoryCache()
	order := fakeOrder(t)
	key := fakeKey(t)
	gotOrder, err := c.Get(key)
	require.Error(t, err)
	require.Nil(t, gotOrder)
	c.cache[key] = &item{
		value:     order,
		createdAt: time.Now().Unix(),
		ttl:       61,
	}
	gotOrder, err = c.Get(key)
	require.NoError(t, err)
	require.Equal(t, order, gotOrder)
}

func TestMemoryCache_Set(t *testing.T) {
	c := NewMemoryCache()
	order := fakeOrder(t)
	key := fakeKey(t)
	require.NoError(t, c.Set(key, order, 1))
	gotOrder, err := c.Get(key)
	require.NoError(t, err)
	require.Equal(t, order, gotOrder)
}

func TestMemoryCache_Clean(t *testing.T) {
	c := NewMemoryCache()
	order := fakeOrder(t)
	key := fakeKey(t)
	require.NoError(t, c.Set(key, order, 1))
	gotOrder, err := c.Get(key)
	require.NoError(t, err)
	require.Equal(t, order, gotOrder)
	time.Sleep(2 * time.Second)
	gotOrder, err = c.Get(key)
	require.Error(t, err)
	require.Nil(t, gotOrder)
}
