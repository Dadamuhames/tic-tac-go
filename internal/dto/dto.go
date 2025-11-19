package dto

type MoveMessage struct {
	MessageType int    `json:messageType`
	FromUserId  int    `json:fromUserId`
	ToUserId    int    `json:toUserId`
	Col         int    `json:col`
	Row         int    `json:row`
	Xo          string `json:xo`
}

type StartRequest struct {
	FromUserId int `json:"fromUserId"`
	ToUserId   int `json:"toUserId"`
}

type StartResponse struct {
	MessageType int    `json:"messageType"`
	FromUserId  int    `json:"fromUserId"`
	ToUserId    int    `json:"toUserId"`
	Xo          string `json:"xo"`
}

type ClientResponse struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type RegisterResponse struct {
	MessageType int              `json:"messageType"`
	User        ClientResponse   `json:"user"`
	Users       []ClientResponse `json:"users"`
}

type WinnerMessage struct {
	MessageType int    `json:"messageType"`
	Winner      string `json:"winner"`
	GameId      int    `json:"gameId"`
}

type CloseMessage struct {
	MessageType int `json:"messageType"`
	UserId      int `json:"userId"`
}
