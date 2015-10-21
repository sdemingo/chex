package tests


import (
	"errors"
	"app/users"
)


const(
	ERR_BADANSWERMARK = "Mark not valid"
)




type Answer struct{
	Id       int64           `json:",string" datastore:"-"`
	QuestId  int64           `json:",string"`
	Quest    *Question
	AuthorId int64           `json:",string"`
	Author   *users.NUser

	Solutions []int
	Comment   string   //for non-test question
	// Timestamp??
}



func NewAnswer(nSolutions int)(Answer){
	a:=new (Answer)
	a.Id = -1
	a.QuestId = -1
	a.AuthorId = -1
	a.Solutions = make([]int,0)
	a.Comment = ""

	return *a
}



func (a Answer) Mark(n int)(error){
	if (a.Quest!=nil) && (n >= len(a.Quest.Options)){
		return errors.New(ERR_BADANSWERMARK)
	}
	a.Solutions = append(a.Solutions,n)
	return nil
}


func (a Answer) IsCorrect (master *Answer)(bool){
	// Compara la respuestas 
	return true
}