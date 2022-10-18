package db

var ErrorUserAlreadyExists = errorUserAlreadyExists{}

type errorUserAlreadyExists struct {
}

func (errorUserAlreadyExists) Error() string {
	return "user already exists"
}
