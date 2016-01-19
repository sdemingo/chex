package tests

import (
	"fmt"

	"model/questions"
	"model/users"

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

// Delete the exercises list of the test
func deleteExercises(wr srv.WrapperRequest, t *Test, err error) error {
	if err != nil {
		return err
	}

	ex := NewExerciseBuffer()
	q := data.NewConn(wr, "tests-exercises")
	q.AddFilter("TestId =", t.Id)
	q.GetMany(&ex)

	for _, e := range ex {
		err := q.Delete(e)
		if err != nil {
			return err
		}
	}

	return nil
}

// Check if the exercises are valid to be in the test. If not, return
// an error to avoid commit them in the database
func checkExercises(wr srv.WrapperRequest, t *Test, err error) error {
	if err != nil {
		return err
	}
	for i := range t.Exercises {
		q, _ := questions.GetQuestById(wr, t.Exercises[i].QuestId)
		if q.Solution.Body.IsUnsolved() {
			return fmt.Errorf("%s", ERR_TESTEXERCISESNOTVALID)
		}
	}
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

	// check if the test will be loaded for a student, in this case
	// do not read the solutions of the questions
	withSolution := true
	if wr.NU.GetRole() < users.ROLE_TEACHER {
		withSolution = false
	}
	// now, load the questions struct for each exercise
	for i := range t.Exercises {
		var q *questions.Question
		if withSolution {
			q, _ = questions.GetQuestById(wr, t.Exercises[i].QuestId)
		} else {
			q, _ = questions.GetQuestByIdWithoutSol(wr, t.Exercises[i].QuestId)
		}
		t.Exercises[i].Quest = q
	}

	return nil
}

func getExerciseById(wr srv.WrapperRequest, ExerciseId int64) (*Exercise, error) {
	return nil, nil
}
