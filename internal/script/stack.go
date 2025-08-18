package script

type Stack struct {
	data []string
}

func NewStack() *Stack {
	return &Stack{make([]string, 0)}
}

func (s *Stack) Push(value string) {
	s.data = append(s.data, value)
}

func (s *Stack) Pop() (value string) {
	value = s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-2]
	return
}

func (s Stack) Size() int {
	return len(s.data)
}
