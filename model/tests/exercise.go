package tests

import (
	"model/questions"

	"appengine/data"
	"appengine/srv"
)

type Exercise struct {
	Id        int64 `datastore:"-"`
	TestId    int64
	QuestId   int64
	Quest     *questions.Question `json:","datastore:"-"`
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

// Add the exercises in the test to the database
func addExercises(wr srv.WrapperRequest, t *Test, err error) error {
	if err != nil {
		return err
	}
	q := data.NewConn(wr, "tests-exercises")
	for _, ex := range t.Exercises {
		ex.TestId = t.Id
		err := q.Put(ex)
		if err != nil {
			return err
		}
	}
	return nil
}

func deleteExercises(wr srv.WrapperRequest, t *Test, err error) error {
	if err != nil {
		return err
	}

	// TODO

	return nil
}

// Fill the exercises array in the test
func loadExercises(wr srv.WrapperRequest, t *Test, err error) error {
	if err != nil {
		return err
	}
	testEx := NewExerciseBuffer()

	qry := data.NewConn(wr, "tests-exercises")
	qry.AddFilter("TestId =", t.Id)
	err = qry.GetMany(&testEx)
	if err != nil {
		return err
	}

	t.Exercises = testEx

	// now, load the questions struct for each exercise
	for i := range t.Exercises {
		q, _ := questions.GetQuestById(wr, t.Exercises[i].QuestId)
		t.Exercises[i].Quest = q
	}

	return nil
}
