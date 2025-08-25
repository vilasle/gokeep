package model

type UserAdd struct {
	Login    string
	Password string
}

type UserUpdate struct {
	ID       int
	Login    string
	Password string
}

type UserInfo struct {
	ID       int
	Login    string
	Password string
}

type PrivateDataSave struct {
	ID     int
	UserID int
	Type   int
	View   string
	Data   []byte
	DEK    []byte
}

type PrivateDataInfo struct {
	ID     int
	UserID int
	Type   int
	View   string
	Data   []byte
	DEK    []byte
}
