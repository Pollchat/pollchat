package main

type Question struct {
	Id        int                  `json:"id"`
	Question  string               `json:"question"`
	Responses map[string]*Response `json:"responses"`
	updated   chan bool
}

type Response struct {
	Response string `json:"response"`
	Count    int    `json:"count"`
}

func GenerateQuestion(id int, q string, r1 string, r2 string, r3 string, r4 string) *Question {
	ques := new(Question)
	ques.Id = id
	ques.Question = q
	ques.Responses = make(map[string]*Response)
	ques.Responses["1"] = &Response{r1, 0}
	ques.Responses["2"] = &Response{r2, 0}
	ques.Responses["3"] = &Response{r3, 0}
	ques.Responses["4"] = &Response{r4, 0}
	return ques
}
