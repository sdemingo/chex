package data

import (
	"fmt"

	"appengine/memcache"
)

func (op *DataConn) PutCacheItem(item DataItem) error {
	citem := new(memcache.Item)
	citem.Key = fmt.Sprintf("%s:%d", op.Entity, item.ID())
	citem.Object = item

	err := memcache.Gob.Set(op.Wreq.C, citem)
	return err
}

func (op *DataConn) GetCacheItem(item DataItem) error {
	key := fmt.Sprintf("%s:%d", op.Entity, item.ID())
	_, err := memcache.Gob.Get(op.Wreq.C, key, item)
	return err
}
