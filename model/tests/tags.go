package tests

import (
	"fmt"

	"model/users"

	"appengine/data"
	"appengine/srv"
)

type TestTag struct {
	Id     int64 `json:",string" datastore:"-"`
	TestId int64
	Tag    string
}

func (t TestTag) ID() int64 {
	return t.Id
}

func (t *TestTag) SetID(id int64) {
	t.Id = id
}

type TestTagBuffer []*TestTag

func NewTestTagBuffer() TestTagBuffer {
	return make([]*TestTag, 0)
}

func (v TestTagBuffer) At(i int) data.DataItem {
	return data.DataItem(v[i])
}

func (v TestTagBuffer) Set(i int, t data.DataItem) {
	v[i] = t.(*TestTag)
}

func (v TestTagBuffer) Len() int {
	return len(v)
}

func addTestTags(wr srv.WrapperRequest, t *Test) error {
	q := data.NewConn(wr, "tests-tags")
	for _, tag := range t.Tags {
		tt := &TestTag{TestId: t.Id, Tag: tag}
		err := q.Put(tt)
		if err != nil {
			return err
		}
	}
	return nil
}

func getTestTags(wr srv.WrapperRequest, t *Test) error {
	var tags []string
	testTags := NewTestTagBuffer()

	qry := data.NewConn(wr, "tests-tags")
	qry.AddFilter("TestId =", t.Id)
	err := qry.GetMany(&testTags)
	if err != nil {
		return err
	}

	tags = make([]string, 0)
	for _, ttag := range testTags {
		tags = append(tags, ttag.Tag)
	}

	t.Tags = tags

	return nil
}

func getAllTestsTags(wr srv.WrapperRequest) ([]string, error) {
	var tagsMap = make(map[string]int, 0)
	var tags = make([]string, 0)
	testTags := NewTestTagBuffer()

	qry := data.NewConn(wr, "tests-tags")
	qry.GetMany(&testTags)

	tags = make([]string, len(testTags))
	for _, qtag := range testTags {
		if _, ok := tagsMap[qtag.Tag]; !ok {
			tagsMap[qtag.Tag] = 1
			tags = append(tags, qtag.Tag)
		}
	}

	return tags, nil
}

func getTestsByTags(wr srv.WrapperRequest, tags []string) (TestBuffer, error) {
	ts := NewTestBuffer()
	ttagsAll := NewTestTagBuffer()

	qry := data.NewConn(wr, "tests-tags")
	qry.GetMany(&ttagsAll)

	filtered := make(map[int64]int)
	for _, tag := range tags {
		for _, tt := range ttagsAll {
			if tt.Tag == tag {
				if _, ok := filtered[tt.TestId]; !ok {
					filtered[tt.TestId] = 1
				} else {
					filtered[tt.TestId]++
				}
			}
		}
	}

	for id, _ := range filtered {
		if filtered[id] == len(tags) {
			q, err := getTestById(wr, fmt.Sprintf("%d", id))
			if err != nil {
				return ts, err
			}

			// only append the questions of the session user
			//if wr.NU.IsAdmin() || q.AuthorId == wr.NU.Id {
			if wr.NU.GetRole() == users.ROLE_ADMIN || q.AuthorId == wr.NU.ID() {
				ts = append(ts, q)
			}
		}
	}

	return ts, nil
}
