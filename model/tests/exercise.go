package tests

import (
	"model/questions"

	"appengine/data"
	"appengine/srv"
)

type Exercise struct {
	Id        int64 `datastore:"-"`
	QuestId   int64
	Quest     questions.Question `json:","datastore:"-"`
	BadPoint  float32
	GoodPoint float32
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
	q := data.NewConn(wr, "tests-exercises")
	for _, ex := range t.Exercises {
		err := q.Put(ex)
		if err != nil {
			return err
		}
	}
	return nil
}
