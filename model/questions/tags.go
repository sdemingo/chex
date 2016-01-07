package questions

import (
	"model/users"

	"appengine/data"
	"appengine/srv"
)

// Question Tag Model
type QuestionTag struct {
	Id      int64 `json:",string" datastore:"-"`
	QuestId int64
	Tag     string
}

func (ut QuestionTag) ID() int64 {
	return ut.Id
}

func (ut *QuestionTag) SetID(id int64) {
	ut.Id = id
}

type QuestionTagBuffer []*QuestionTag

func NewQuestionTagBuffer() QuestionTagBuffer {
	return make([]*QuestionTag, 0)
}

func (v QuestionTagBuffer) At(i int) data.DataItem {
	return data.DataItem(v[i])
}

func (v QuestionTagBuffer) Set(i int, t data.DataItem) {
	v[i] = t.(*QuestionTag)
}

func (v QuestionTagBuffer) Len() int {
	return len(v)
}

// Return the questions that have all tags in the array. The questions
// must be wrote by  the session user
func getQuestByTags(wr srv.WrapperRequest, tags []string) (QuestionBuffer, error) {
	qs := NewQuestionBuffer()
	qTagsAll := NewQuestionTagBuffer()

	qry := data.NewConn(wr, "questions-tags")
	qry.GetMany(&qTagsAll)

	filtered := make(map[int64]int)
	for _, tag := range tags {
		for _, qt := range qTagsAll {
			if qt.Tag == tag {
				if _, ok := filtered[qt.QuestId]; !ok {
					filtered[qt.QuestId] = 1
				} else {
					filtered[qt.QuestId]++
				}
			}
		}
	}

	for id, _ := range filtered {
		if filtered[id] == len(tags) {
			q, err := GetQuestById(wr, id)
			if err != nil {
				return qs, err
			}

			// only append the questions of the session user
			//if wr.NU.IsAdmin() || q.AuthorId == wr.NU.Id {
			if wr.NU.GetRole() == users.ROLE_ADMIN || q.AuthorId == wr.NU.ID() {
				qs = append(qs, q)
			}
		}
	}

	return qs, nil
}

// Add the tags of the question to the database
func addQuestTags(wr srv.WrapperRequest, q *Question) error {
	for _, tag := range q.Tags {
		qry := data.NewConn(wr, "questions-tags")
		qt := &QuestionTag{QuestId: q.Id, Tag: tag}
		err := qry.Put(qt)
		if err != nil {
			return err
		}
	}
	return nil
}

// Return the tags of the question
func getQuestTags(wr srv.WrapperRequest, q *Question) ([]string, error) {
	var tags []string
	questionTags := NewQuestionTagBuffer()

	qry := data.NewConn(wr, "questions-tags")
	qry.AddFilter("QuestId =", q.Id)
	err := qry.GetMany(&questionTags)
	if err != nil {
		return tags, err
	}

	tags = make([]string, 0)
	for _, qtag := range questionTags {
		tags = append(tags, qtag.Tag)
	}

	return tags, nil

}

// Return all questions tags of all questions in the database
func getAllQuestionsTags(wr srv.WrapperRequest) ([]string, error) {
	var tagsMap = make(map[string]int, 0)
	var tags = make([]string, 0)
	questionTags := NewQuestionTagBuffer()

	qry := data.NewConn(wr, "questions-tags")
	qry.GetMany(&questionTags)

	tags = make([]string, len(questionTags))
	for _, qtag := range questionTags {
		if _, ok := tagsMap[qtag.Tag]; !ok {
			tagsMap[qtag.Tag] = 1
			tags = append(tags, qtag.Tag)
		}
	}

	return tags, nil
}

// Return all questions tags from all questions of the author with authorId
func getQuestionsTagsFromUser(wr srv.WrapperRequest, authorId int64) ([]string, error) {
	var tagsMap = make(map[string]int, 0)
	userQuests, err := getQuestByAuthor(wr, authorId)

	tags := make([]string, 0)
	for _, q := range userQuests {
		for _, tag := range q.Tags {
			if _, ok := tagsMap[tag]; !ok {
				tagsMap[tag] = 1
				tags = append(tags, tag)
			}
		}
	}

	return tags, err
}
