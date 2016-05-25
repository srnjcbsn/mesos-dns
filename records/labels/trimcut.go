package labels

import (
	"bytes"
)

var allowedChars = []byte("0123456789abcdefghijklmnopqrstuvwxyz")

type state struct {
	name, accum []byte
	left        []byte
	right       string
	maxlen      int
	pos         int
}

type stateFn func(*state) stateFn

func newState(name []byte, maxlen int, left, right string) *state {
	return &state{
		name:   bytes.ToLower(name),
		accum:  make([]byte, 0, len(name)), // better cap: min(maxlen, len(name))
		maxlen: maxlen,
		left:   []byte(left),
		right:  right,
	}
}

func (s *state) run() string {
	f := initialState
	for f != nil {
		f = f(s)
	}
	return string(s.accum)
}

func (s *state) trimRight(x string, f stateFn) stateFn {
	s.accum = bytes.TrimRight(s.accum, x)
	return f
}

func initialState(s *state) stateFn {
	if s.pos >= len(s.name) {
		return s.trimRight(s.right, nil)
	} else if bytes.IndexByte(allowedChars, s.name[s.pos]) > -1 && bytes.IndexByte(s.left, s.name[s.pos]) == -1 {
		return middleState
	}
	s.pos++
	return initialState
}

func middleState(s *state) stateFn {
	if s.pos >= len(s.name) {
		return s.trimRight(s.right, nil)
	} else if len(s.accum) >= s.maxlen {
		return s.trimRight("-", endState)
	} else if s.name[s.pos] == '-' || s.name[s.pos] == '_' || s.name[s.pos] == '.' {
		s.accum = append(s.accum, '-')
	} else if bytes.IndexByte(allowedChars, s.name[s.pos]) > -1 {
		s.accum = append(s.accum, s.name[s.pos])
	}
	s.pos++
	return middleState
}

func endState(s *state) stateFn {
	if s.pos >= len(s.name) || len(s.accum) == s.maxlen {
		return s.trimRight(s.right, nil)
	} else if bytes.IndexByte(allowedChars, s.name[s.pos]) > -1 {
		s.accum = append(s.accum, s.name[s.pos])
	}
	s.pos++
	return endState
}
