package dice

import (
	"fmt"
	"math/rand"
)

type DiceExpression interface {
	DisplayString() string
	Eval() int
	/* is this statement (and all its children) "literal"? Literal numeric
	 * expressions and compound expressions whose children are both literals are
	 * literal, diceExpressions are not
	 */
	Literal() bool
}

/*
diceExpression is an expression for dicerolls
(MdN expressions, e.g. 2d4)
*/
type diceExpression struct {
	count int
	sides int
}

func NewDiceExpression(count int, sides int) DiceExpression {
	return &diceExpression{
		count: count,
		sides: sides,
	}
}

func (exp *diceExpression) DisplayString() string {
	return fmt.Sprintf("%dd%d", exp.count, exp.sides)
}

func (exp *diceExpression) Eval() int {
	total := 0
	for i := 0; i < exp.count; i++ {
		total += rand.Intn(exp.sides) + 1
	}
	return total
}

func (exp *diceExpression) Literal() bool {
	return false
}

/*
literalExpression is an expression for literal numbers
*/
type literalExpression struct {
	val int
}

func NewLiteralExpression(val int) DiceExpression {
	return &literalExpression{val: val}
}

func (exp *literalExpression) DisplayString() string {
	return fmt.Sprintf("%d", exp.val)
}

func (exp *literalExpression) Eval() int {
	return exp.val
}

func (exp *literalExpression) Literal() bool {
	return true
}

/*
additionExpression is an expression that expresses addition
*/
type additionExpression struct {
	lhs DiceExpression
	rhs DiceExpression
}

func NewAdditionExpression(lhs DiceExpression, rhs DiceExpression) DiceExpression {
	return &additionExpression{
		lhs: lhs,
		rhs: rhs,
	}
}

func (exp *additionExpression) DisplayString() string {
	return fmt.Sprintf("%v + %v", exp.lhs.DisplayString(), exp.rhs.DisplayString())
}

func (exp *additionExpression) Eval() int {
	return exp.lhs.Eval() + exp.rhs.Eval()
}

func (exp *additionExpression) Literal() bool {
	return exp.lhs.Literal() && exp.rhs.Literal()
}
