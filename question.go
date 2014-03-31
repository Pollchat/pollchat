package main

type Question struct {
	Id        int               `json:"id"`
	Question  string            `json:"question"`
	Responses map[int]*Response `json:"responses"`
	updated   chan bool
}

type Response struct {
	Response string `json:"response"`
	Count    int    `json:"count"`
}

var questions = make(map[int]*Question)
var qcount = 0

func GenerateQuestion(q string, r1 string, r2 string, r3 string, r4 string) *Question {
	ques := new(Question)
	qcount++
	ques.Id = qcount
	ques.Question = q
	ques.Responses = make(map[int]*Response)
	ques.Responses[1] = &Response{r1, 0}
	ques.Responses[2] = &Response{r2, 0}
	ques.Responses[3] = &Response{r3, 0}
	ques.Responses[4] = &Response{r4, 0}
	return ques
}
