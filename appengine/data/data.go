package data

import (
	"appengine/datastore"
	"appengine/srv"
)

type DataItem interface {
	ID() int64
	SetID(id int64)
}

type BufferItems interface {
	At(i int) DataItem
	Set(i int, d DataItem)
	Len() int
}

type DataConn struct {
	Entity string
	Query  *datastore.Query
	Wreq   srv.WrapperRequest
}

func NewConn(wr srv.WrapperRequest, entity string) *DataConn {
	op := new(DataConn)
	op.Entity = entity
	op.Wreq = wr
	op.Query = datastore.NewQuery(entity)
	return op
}

func (op *DataConn) AddFilter(filter string, value interface{}) {
	op.Query = op.Query.Filter(filter, value)
}

func (op *DataConn) Put(obj DataItem) error {

	var key *datastore.Key
	c := op.Wreq.C

	if id := obj.ID(); id > 0 {
		key = datastore.NewKey(c, op.Entity, "", obj.ID(), nil)
	} else {
		key = datastore.NewIncompleteKey(c, op.Entity, nil)
	}

	key, err := datastore.Put(c, key, obj)
	obj.SetID(key.IntID())

	return err
}

func (op *DataConn) GetMany(items BufferItems) error {
	c := op.Wreq.C
	keys, err := op.Query.GetAll(c, items)
	if err != nil {
		return err
	}

	for i := range keys {
		it := items.At(i)
		it.SetID(keys[i].IntID())
	}

	return err
}

func (op *DataConn) Get(item DataItem) error {
	c := op.Wreq.C
	key := datastore.NewKey(c, op.Entity, "", item.ID(), nil)
	err := datastore.Get(c, key, item)
	item.SetID(key.IntID())
	return err
}

func (op *DataConn) Delete(item DataItem) error {
	c := op.Wreq.C
	key := datastore.NewKey(c, op.Entity, "", item.ID(), nil)
	return datastore.Delete(c, key)
}
