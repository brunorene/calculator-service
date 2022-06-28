package operator

type Operator interface {
	Result(left, right int) int
}

type Add struct{}

func (op *Add) Result(left, right int) int {
	return left + right
}

type Subtract struct{}

func (op *Subtract) Result(left, right int) int {
	return left - right
}

type Multiply struct{}

func (op *Multiply) Result(left, right int) int {
	return left * right
}

type Divide struct{}

func (op *Divide) Result(left, right int) int {
	return left / right
}
