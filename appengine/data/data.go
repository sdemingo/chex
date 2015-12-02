package data

import (
	"appengine"
	"appengine/datastore"
)

type DataItem interface {
	ID() int64
	SetID(id int64)
	Kind() string
}

type VectorItems interface {
	At(i int) DataItem
	Set(i int, d DataItem)
	Len() int
}

type DataConn struct {
	Entity string
	Query  *datastore.Query
	Ids    []int64
	Items  []DataItem
}

func (op *DataConn) AddFilter(filter string, value interface{}) {
	op.Query = op.Query.Filter(filter, value)
}

func (op *DataConn) Put(c appengine.Context, obj DataItem) (int64, error) {

	var key *datastore.Key

	if id := obj.ID(); id != 0 {
		key = datastore.NewKey(c, op.Entity, "", obj.ID(), nil)
	} else {
		key = datastore.NewIncompleteKey(c, op.Entity, nil)
	}

	key, err := datastore.Put(c, key, obj)
	obj.SetID(key.IntID())

	return key.IntID(), err
}

func (op *DataConn) GetMany(c appengine.Context, items interface{}) error {

	keys, err := op.Query.GetAll(c, items)
	if err != nil {
		return err
	}

	for i := range keys {
		op.Ids = append(op.Ids, keys[i].IntID())
	}

	return err
}

func (op *DataConn) Get(c appengine.Context, item DataItem) error {
	key := datastore.NewKey(c, op.Entity, "", item.ID(), nil)
	return datastore.Get(c, key, item)
}
