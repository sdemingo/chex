package questions

import (
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
	"html/template"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/russross/blackfriday"

	"app/users"
	"appengine/datastore"
	"appengine/srv"
)

const (
	ERR_NOTVALIDQUEST     = "Pregunta no valido"
	ERR_DUPLICATEDQUEST   = "Pregunta duplicada"
	ERR_QUESTNOTFOUND     = "Pregunta no encontrada"
	ERR_BADRENDEREDANSWER = "Pregunta renderizada erroneamente"

	TMPL_NOANSWEREDQUESTION = `
		<ul>{{range .}}<li>{{ . }}</li>{{end}}</ul>
`
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

func (q *Question) GetHTMLOptions() template.HTML {
	var s template.HTML
	if q.Solution != nil {
		// llamamos al renderizador de la respuesta
	} else {
		var doc bytes.Buffer
		t, err := template.New("options").Parse(TMPL_NOANSWEREDQUESTION)
		err = t.Execute(&doc, q.Options)
		if err != nil {
			s = template.HTML(ERR_BADRENDEREDANSWER)
		} else {
			s = template.HTML(doc.String())
		}
	}
	return s
}

func getQuestions(wr srv.WrapperRequest, filters map[string][]string) ([]Question, error) {

	var qs []Question
	var err error

	if filters["id"] != nil {
		q, err := getQuestById(wr, filters["id"][0])
		qs := make([]Question, 1)
		qs[0] = q
		return qs, err
	}

	if filters["tags"] != nil {
		qs, err := getQuestByTags(wr, strings.Split(filters["tags"][0], ","))
		return qs, err
	}

	return qs, err
}

func putQuestion(wr srv.WrapperRequest, q Question) error {
	if err := q.IsValid(); err != nil {
		return err
	}

	q.TimeStamp = time.Now()
	q.SetCheckSum()
	q.AuthorId = wr.NU.Id

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

func getQuestById(wr srv.WrapperRequest, s_id string) (Question, error) {

	var q Question

	id, err := strconv.ParseInt(s_id, 10, 64)
	if err != nil {
		return q, errors.New(ERR_QUESTNOTFOUND)
	}

	if id != 0 {
		k := datastore.NewKey(wr.C, "questions", "", id, nil)
		datastore.Get(wr.C, k, &q)
	} else {
		return q, errors.New(ERR_QUESTNOTFOUND)
	}

	q.Id = id
	q.Tags, _ = getQuestTags(wr, q)

	return q, nil
}

func getQuestByTags(wr srv.WrapperRequest, tags []string) ([]Question, error) {
	var qs []Question
	var qTagsAll []QuestionTag

	qry := datastore.NewQuery("questions-tags")
	_, err := qry.GetAll(wr.C, &qTagsAll)
	if err != nil {
		return qs, errors.New(ERR_QUESTNOTFOUND)
	}

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
			if q.AuthorId == wr.NU.Id {
				qs = append(qs, q)
			}
		}
	}

	return qs, nil
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

func getAllQuestionsTags(wr srv.WrapperRequest) ([]string, error) {
	var tagsMap = make(map[string]int, 0)
	var questionTags []QuestionTag
	var tags = make([]string, 0)

	q := datastore.NewQuery("questions-tags")
	_, err := q.GetAll(wr.C, &questionTags)
	if err != nil {
		return tags, err
	}
	tags = make([]string, len(questionTags))
	for _, qtag := range questionTags {
		if _, ok := tagsMap[qtag.Tag]; !ok {
			tagsMap[qtag.Tag] = 1
			tags = append(tags, qtag.Tag)
		}
	}

	return tags, nil
}
