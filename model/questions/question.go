package questions

import (
	"crypto/md5"
	"errors"
	"fmt"
	"html/template"
	"regexp"
	"strconv"
	"strings"
	"time"

	"app/users"
	"model/answers"

	"github.com/russross/blackfriday"

	"appengine/data"
	"appengine/srv"
)

const (
	ERR_NOTVALIDQUEST   = "Pregunta no valido"
	ERR_DUPLICATEDQUEST = "Pregunta duplicada"
	ERR_QUESTNOTFOUND   = "Pregunta no encontrada"
)

// Question Model

type Question struct {
	Id         int64           `json:",string" datastore:"-"`
	AuthorId   int64           `json:",string"`
	Author     *users.NUser    `datastore:"-"`
	SolutionId int64           `json:",string"`
	Solution   *answers.Answer `datastore:"-"`
	TimeStamp  time.Time       `json:"`
	CheckSum   string          `json:"`

	AType   answers.AnswerBodyType `json:",string"`
	Text    string
	Hint    string
	Options []string
	Tags    []string
}

func NewQuestion() *Question {
	q := new(Question)
	q.TimeStamp = time.Now()
	q.Options = make([]string, 0)
	q.Tags = make([]string, 0)
	return q
}

func (q *Question) SetSolution(sol *answers.Answer) {
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

func (ut Question) ID() int64 {
	return ut.Id
}

func (ut *Question) SetID(id int64) {
	ut.Id = id
}

func (q *Question) SetCheckSum() {
	re := regexp.MustCompile("\\s+")

	s := re.ReplaceAllString(strings.ToLower(q.Text), "")
	for _, op := range q.Options {
		s = s + re.ReplaceAllString(strings.ToLower(op), "")
	}
	q.CheckSum = fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

func (q *Question) IsValid() error {
	if q != nil && q.Text == "" {
		return errors.New(ERR_NOTVALIDQUEST)
	}
	return nil
}

func (q *Question) GetHTMLText() template.HTML {
	in := []byte(q.Text)
	return template.HTML(string(blackfriday.MarkdownBasic(in)))
}

type QuestionBuffer []*Question

func NewQuestionBuffer() QuestionBuffer {
	return make([]*Question, 0)
}

func (v QuestionBuffer) At(i int) data.DataItem {
	return data.DataItem(v[i])
}

func (v QuestionBuffer) Set(i int, t data.DataItem) {
	v[i] = t.(*Question)
}

func (v QuestionBuffer) Len() int {
	return len(v)
}

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

// Data backend access functions

func getQuestions(wr srv.WrapperRequest, filters map[string][]string) (QuestionBuffer, error) {
	qs := NewQuestionBuffer()
	var err error

	if filters["id"] != nil {
		q, err := getQuestById(wr, filters["id"][0])
		qs = append(qs, q)
		return qs, err
	}

	if filters["tags"] != nil {
		qs, err := getQuestByTags(wr, strings.Split(filters["tags"][0], ","))
		return qs, err
	}

	return qs, err
}

func putQuestion(wr srv.WrapperRequest, q *Question) error {
	if err := q.IsValid(); err != nil {
		return err
	}

	q.TimeStamp = time.Now()
	q.SetCheckSum()
	q.AuthorId = wr.NU.Id

	_, err := getQuestByChecksum(wr, q.CheckSum)
	if err == nil {
		return errors.New(ERR_DUPLICATEDQUEST)
	}

	qc := data.NewConn(wr, "questions")
	qc.Put(q)

	// Add a QuestionsTags entry for each tag for this questions
	err = addQuestTags(wr, q)

	return err
}

func getQuestByChecksum(wr srv.WrapperRequest, sum string) (*Question, error) {
	qs := NewQuestionBuffer()
	q := NewQuestion()

	qc := data.NewConn(wr, "questions")
	qc.AddFilter("CheckSum =", sum)
	err := qc.GetMany(&qs)
	if err != nil {
		return nil, err
	}
	if len(qs) == 0 {
		return nil, errors.New(ERR_QUESTNOTFOUND)
	}
	q = qs[0]
	q.Tags, _ = getQuestTags(wr, q)

	return q, nil
}

func getQuestById(wr srv.WrapperRequest, s_id string) (*Question, error) {
	q := NewQuestion()
	var err error

	q.Id, err = strconv.ParseInt(s_id, 10, 64)
	if err != nil {
		return q, errors.New(ERR_QUESTNOTFOUND)
	}

	qry := data.NewConn(wr, "questions")
	qry.Get(q)

	q.Tags, _ = getQuestTags(wr, q)

	// search the solution. An answer for this quest from the same author
	var errA error
	q.Solution, errA = answers.GetAnswer(wr, q.AuthorId, q.Id)
	if errA == nil {
		q.SolutionId = q.Solution.Id
	} else {
		q.SolutionId = -1
	}

	return q, err
}

func getQuestByAuthor(wr srv.WrapperRequest, authorId string) (QuestionBuffer, error) {
	qs := NewQuestionBuffer()

	id, err := strconv.ParseInt(authorId, 10, 64)
	if err != nil {
		return qs, errors.New(ERR_QUESTNOTFOUND)
	}

	qry := data.NewConn(wr, "questions")
	qry.AddFilter("AuthorId =", id)
	qry.GetMany(&qs)

	return qs, nil
}

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
			q, err := getQuestById(wr, fmt.Sprintf("%d", id))
			if err != nil {
				return qs, err
			}

			// only append the questions of the session user
			if wr.NU.IsAdmin() || q.AuthorId == wr.NU.Id {
				qs = append(qs, q)
			}
		}
	}

	return qs, nil
}

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

func getQuestionsTagsFromUser(wr srv.WrapperRequest, authorId int64) ([]string, error) {
	var tagsMap = make(map[string]int, 0)
	userQuests, err := getQuestByAuthor(wr, fmt.Sprintf("%d", authorId))

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
