package tags

import (
	//"fmt"

	"appengine/datastore"
	"appengine/srv"
)

type Tag struct {
	Id   int64 `json:",string" datastore:"-"`
	Name string
}

func getAllTags(wr srv.WrapperRequest) (tags []Tag, err error) {

	return
}

func putTag(wr srv.WrapperRequest, tag Tag) error {

	key := datastore.NewKey(wr.C, "tags", "", 0, nil)
	key, err := datastore.Put(wr.C, key, &tag)
	if err != nil {
		return err
	}
	tag.Id = key.IntID()

	return nil
}

func updateTag(wr srv.WrapperRequest, tag Tag) error {

	return putTag(wr, tag)
}

func deleteTag(wr srv.WrapperRequest, tag Tag) error {
	key := datastore.NewKey(wr.C, "tags", "", tag.Id, nil)
	return datastore.Delete(wr.C, key)
}
