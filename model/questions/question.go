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

	"model/answers"
	"model/users"

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
	Id        int64           `json:",string" datastore:"-"`
	Author    *users.NUser    `datastore:"-"`
	Solution  *answers.Answer `datastore:"-"`
	AuthorId  int64           `json:",string"`
	TimeStamp time.Time       `json:"`
	CheckSum  string          `json:"`

	AType   answers.AnswerBodyType `json:",string"`
	Text    string
	Hint    string
	Options []string
	Tags    []string `datastore:"-"`
}

func NewQuestion() *Question {
	q := new(Question)
	q.TimeStamp = time.Now()
	q.Options = make([]string, 0)
	q.Tags = make([]string, 0)
	return q
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

func (q *Question) GetHTMLSolution() template.HTML {
	if q.Solution != nil {
		_, solved, _ := q.Solution.Body.GetHTML(q.Options)
		return solved
	}
	return template.HTML("")
}

func (q *Question) GetHTMLAnswer() template.HTML {
	if q.Solution != nil {
		unSolved, _, _ := q.Solution.Body.GetHTML(q.Options)
		return unSolved
	}
	return template.HTML("")
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

func getQuestions(wr srv.WrapperRequest, filters map[string][]string) (QuestionBuffer, error) {
	qs := NewQuestionBuffer()
	var err error

	if filters["id"] != nil {
		id, err := strconv.ParseInt(filters["id"][0], 10, 64)
		if err != nil {
			return qs, errors.New(ERR_NOTVALIDQUEST)
		}
		q, err := GetQuestById(wr, id)
		qs = append(qs, q)
		return qs, err
	}

	if filters["tags"] != nil {
		qs, err := getQuestByTags(wr, strings.Split(filters["tags"][0], ","))
		return qs, err
	}

	return qs, err
}

// Return a question with the id
func GetQuestById(wr srv.WrapperRequest, id int64) (*Question, error) {
	q := NewQuestion()
	var err error

	q.Id = id

	qry := data.NewConn(wr, "questions")
	qry.Get(q)

	q.Tags, _ = getQuestTags(wr, q)

	// search the solution. An answer for this quest from the same
	// author. If not exits we create a dummy answer body of
	// questions Atype. A dummy answer body always must return false
	// in its UnSolved method
	q.Solution, _ = answers.GetSolutionAnswer(wr, q.AuthorId, q.Id)
	if q.Solution == nil {
		q.Solution, err = answers.NewAnswerWithBody(-1, -1, q.AType)
	}

	return q, err
}

func GetQuestByIdWithoutSol(wr srv.WrapperRequest, id int64) (*Question, error) {
	q := NewQuestion()
	var err error

	q.Id = id
	qry := data.NewConn(wr, "questions")
	qry.Get(q)
	q.Tags, _ = getQuestTags(wr, q)
	q.Solution, err = answers.NewAnswerWithBody(-1, -1, q.AType)
	return q, err
}

// Write a new question on the database
func putQuestion(wr srv.WrapperRequest, q *Question) error {
	if err := q.IsValid(); err != nil {
		return err
	}

	q.TimeStamp = time.Now()
	q.SetCheckSum()
	q.AuthorId = wr.NU.ID()

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

// Return a question with the checksum
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

// Return all questions of the author with authorId
func getQuestByAuthor(wr srv.WrapperRequest, authorId int64) (QuestionBuffer, error) {
	qs := NewQuestionBuffer()

	qry := data.NewConn(wr, "questions")
	qry.AddFilter("AuthorId =", authorId)
	qry.GetMany(&qs)

	for i := range qs {
		qs[i].Tags, _ = getQuestTags(wr, qs[i])
	}

	return qs, nil
}
