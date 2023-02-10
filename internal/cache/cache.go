package cache

import (
	"errors"
	"log"
	"strconv"
	"sync"

	"github.com/afterglowflexin/wildberries/level0/internal/store"
	"github.com/jackc/pgx/v5"
)

type Cache struct {
	sync.RWMutex
	items map[string]string
}

func New(conn *pgx.Conn) *Cache {
	items := store.GetAllOrders(conn)
	if items == nil {
		items = make(map[string]string)
	}

	cache := Cache{
		items: items,
	}

	log.Println("Initiated cache")
	return &cache
}

func (c *Cache) AddOrder(order string) {
	id := c.findCurrID()
	log.Println("[DEBUG] adding order in cache")

	c.Lock()

	defer c.Unlock()

	c.items[id] = order
}

func (c *Cache) GetOrder(id string) (string, error) {
	c.RLock()

	defer c.RUnlock()

	order, found := c.items[id]

	if !found {
		return "", errors.New("order not found")
	}

	return order, nil
}

func (c *Cache) findCurrID() string {
	prevID := len(c.items)

	return strconv.Itoa(prevID + 1)
}
