package db

var ErrorUserAlreadyExists = errorUserAlreadyExists{}

type errorUserAlreadyExists struct {
}

func (errorUserAlreadyExists) Error() string {
	return "user already exists"
}

var ErrorUserNotFound = errorUserNotFound{}

type errorUserNotFound struct {
}

func (errorUserNotFound) Error() string {
	return "user not found"
}
