package questions

import (
	"crypto/md5"
	"errors"
	"fmt"
	"time"

	"app/users"
	"appengine/datastore"
	"appengine/srv"
)

const (
	ERR_NOTVALIDQUEST   = "Pregunta no valido"
	ERR_DUPLICATEDQUEST = "Pregunta duplicada"
	ERR_QUESTNOTFOUND   = "Pregunta no encontrada"
)

type QuestionTag struct {
	Id      int64 `json:",string" datastore:"-"`
	QuestId int64
	Tag     string
}

type Question struct {
	Id         int64        `json:",string" datastore:"-"`
	AuthorId   int64        `json:",string"`
	Author     *users.NUser `datastore:"-"`
	SolutionId int64        `json:",string"`
	Solution   *Answer      `datastore:"-"`
	TimeStamp  time.Time    `json:"`
	CheckSum   string       `json:"`

	AType   AnswerBodyType `json:",string"`
	Text    string
	Hint    string
	Options []string
	Tags    []string
}

func NewQuestion(text string, options []string, tags []string) Question {
	q := new(Question)
	q.Id = -1
	q.AuthorId = -1
	q.SolutionId = -1

	q.Text = text
	q.Hint = ""
	q.Options = options
	q.Tags = tags

	return *q
}

func (q *Question) SetSolution(sol *Answer) {
	if sol != nil {
		q.Solution = sol
		q.SolutionId = sol.Id
	}
}

func (q *Question) SetAuthor(author *users.NUser) {
	if author != nil {
		q.Author = author
		q.AuthorId = author.Id
	}
}

func (q *Question) SetCheckSum() {
	s := q.Text
	for _, op := range q.Options {
		s = s + op
	}
	q.CheckSum = fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

func (q *Question) IsValid() error {
	if q != nil && q.Text == "" {
		return errors.New(ERR_NOTVALIDQUEST)
	}
	return nil
}

func putQuestion(wr srv.WrapperRequest, q Question) error {
	if err := q.IsValid(); err != nil {
		return err
	}

	q.TimeStamp = time.Now()
	q.SetCheckSum()

	_, err := getQuestionByChecksum(wr, q.CheckSum)
	if err == nil {
		return errors.New(ERR_DUPLICATEDQUEST)
	}

	key := datastore.NewKey(wr.C, "questions", "", 0, nil)
	key, err = datastore.Put(wr.C, key, &q)
	if err != nil {
		return err
	}
	q.Id = key.IntID()

	// Add a QuestionsTags entry for each tag for this questions
	addQuestTags(wr, q)

	return nil
}

func getQuestionByChecksum(wr srv.WrapperRequest, sum string) (Question, error) {
	var qs []Question
	var q Question

	qry := datastore.NewQuery("questions").Filter("Checksum =", sum)

	keys, err := qry.GetAll(wr.C, &qs)
	if (len(keys) == 0) || err != nil {
		return q, errors.New(ERR_QUESTNOTFOUND)
	}
	q = qs[0]
	q.Id = keys[0].IntID()
	q.Tags, _ = getQuestTags(wr, q)

	return q, nil
}

func addQuestTags(wr srv.WrapperRequest, q Question) error {
	for _, tag := range q.Tags {
		key := datastore.NewKey(wr.C, "questions-tags", "", 0, nil)
		qt := QuestionTag{QuestId: q.Id, Tag: tag}
		key, err := datastore.Put(wr.C, key, &qt)
		if err != nil {
			return err
		}
	}
	return nil
}

func getQuestTags(wr srv.WrapperRequest, q Question) ([]string, error) {

	var tags []string
	var questionTags []QuestionTag

	qry := datastore.NewQuery("questions-tags").Filter("QuestId =", q.Id)
	_, err := qry.GetAll(wr.C, &questionTags)
	if err != nil {
		return tags, err
	}
	tags = make([]string, 0)
	for _, qtag := range questionTags {
		tags = append(tags, qtag.Tag)
	}

	return tags, nil

}
