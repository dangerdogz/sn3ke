package main

type Message struct {
	Owner   string `json:"owner"`
	Content string `json:"content"`
}

func (self *Message) String() string {
	return self.Owner + self.Content
}
