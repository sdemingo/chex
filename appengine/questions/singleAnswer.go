package questions

// Single Test Answer

type SingleTestBody struct {
	Id       int64          `json:",string" datastore:"-"`
	atype    AnswerBodyType `json:",string"`
	Solution int
}

func NewSingleTestBody(sol int) SingleTestBody {
	return SingleTestBody{-1, -1, sol}
}

func (a SingleTestBody) Equals(master SingleTestBody) bool {
	return a.Solution == master.Solution
}
