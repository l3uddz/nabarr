package nabarr

import (
	"time"

	"github.com/antonmedv/expr/vm"
	"github.com/l3uddz/nabarr/media"
)

type ExprProgram struct {
	expression string
	Program    *vm.Program
}

func (p *ExprProgram) String() string {
	return p.expression
}

func NewExprProgram(expression string, vm *vm.Program) *ExprProgram {
	return &ExprProgram{
		expression: expression,
		Program:    vm,
	}
}

type ExprEnv struct {
	media.Item
	Now func() time.Time
}

func NewExprEnv(media *media.Item) *ExprEnv {
	return &ExprEnv{
		Item: *media,
		Now:  func() time.Time { return time.Now().UTC() },
	}
}
