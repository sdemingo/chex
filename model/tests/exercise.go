package tests

import (
	"model/questions"

	"appengine/data"
	"appengine/srv"
)

type Exercise struct {
	Id        int64              `json:",string" datastore:"-"`
	QuestId   int64              `json:",string"`
	Quest     questions.Question `datastore:"-"`
	BadPoint  int
	GoodPoint int
}

func NewExercise() *Exercise {
	e := new(Exercise)
	return e
}

func (e Exercise) ID() int64 {
	return e.Id
}

func (e *Exercise) SetID(id int64) {
	e.Id = id
}

type ExerciseBuffer []*Exercise

func NewExerciseBuffer() ExerciseBuffer {
	return make([]*Exercise, 0)
}

func (v ExerciseBuffer) At(i int) data.DataItem {
	return data.DataItem(v[i])
}

func (v ExerciseBuffer) Set(i int, t data.DataItem) {
	v[i] = t.(*Exercise)
}

func (v ExerciseBuffer) Len() int {
	return len(v)
}

func addExercises(wr srv.WrapperRequest, t *Test) error {
	q := data.NewConn(wr, "exercises")
	for _, ex := range t.Exercises {
		err := q.Put(ex)
		if err != nil {
			return err
		}
	}
	return nil
}
