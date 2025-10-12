package lib

import "testing"

func TestStack(t *testing.T) {
	stack := NewStack[int]()

	if stack.Size() != 0 {
		t.Errorf(`Faild, unmatched sizes:
				  expected: 0,
				  actual: %v`, stack.Size())

		return
	}

	stack.Push(10)

	if stack.Size() != 1 {
		t.Errorf(`Faild, unmatched sizes:
				  expected: 1,
				  actual: %v`, stack.Size())

		return
	}

	val, _ := stack.Peek()
	if val != 10 {
		t.Errorf(`Faild, uncorrect peeked value:
				  expected: 10,
				  actual: %v`, val)

		return
	}

	last, _ := stack.Pop()

	if last != 10 {
		t.Errorf(`Faild, uncorrect poped value:
				  expected: 10,
				  actual: %v`, last)

		return
	}

	_, err := stack.Pop()
	if err == nil {
		t.Error("Faild, expectecd an error")

		return
	}

	_, err = stack.Peek()
	if err == nil {
		t.Error("Faild, expectecd an error")

		return
	}
}

func TestMoveStackToStack(t *testing.T) {
	stackA := NewStack[int]()
	stackA.Push(10)
	stackA.Push(9)
	stackA.Push(8)

	stackB := NewStack[int]()
	stackB.Push(7)

	n, err := MoveStackToStack(&stackA, &stackB)
	if err != nil {
		t.Errorf("Faild, expected no errors got: %v", err)

		return
	}

	if n != 3 {
		t.Errorf(`Faild, unmatch moved items form stackA to stackB
				  expected: 3,
				  actual: %v`, n)

		return
	}

	val, _ := stackB.Peek()
	if val != 10 {
		t.Errorf(`Faild, uncorrect peeked value:
				  expected: 10,
				  actual: %v`, val)

		return
	}
}

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "",
			expected: []string{},
		},
		{
			input:    "how are you",
			expected: []string{"how", "are", "you"},
		},
		{
			input:    "Hello World",
			expected: []string{"hello", "world"},
		},
		{
			input:    "Charmander Bulbasaur PIKACHU",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
		{
			input:    "    I     love     VIDEO games   ",
			expected: []string{"i", "love", "video", "games"},
		},
	}

	for _, c := range cases {
		actual := CleanInput(c.input)

		if len(actual) != len(c.expected) {
			t.Errorf(`
			Faild, unmatched lenghts, got len: %v, expected len: %v
			======================================
			input: %v`, len(actual), len(c.expected), c.input)

			continue
		}

		for i := range actual {
			if actual[i] != c.expected[i] {
				t.Errorf(`
			Faild, unmatched world, got: %v, expected word: %v
			======================================
			input: %v
			actual: %v
			expected: %v`, actual[i], c.expected[i], c.input, actual, c.expected)

				continue
			}
		}
	}
}
