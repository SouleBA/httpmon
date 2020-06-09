package monitor

import (
	"errors"
	"path"
	"strings"
)

type entry struct {
	remoteIP string
	request
	responseSize uint64
	statusCode   string
}

type request struct {
	method   string
	section  string
	protocol string
}

func NewRequest(requestString string) (*request, error) {
	s := strings.Split(requestString, " ")
	if len(s) != 3 {
		return nil, errors.New("Wrong input. Was expecting \"method path protocol \"")
	}

	var section string

	if dir, _ := path.Split(s[1]); dir != "/" {
		section = dir[:len(dir)-1]
	} else {
		section = s[1]
	}

	return &request{
		method:   s[0],
		section:  section,
		protocol: s[2],
	}, nil
}
