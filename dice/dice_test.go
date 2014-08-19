package dice

import (
  "math/rand"
  "time"
  "fmt"
  "testing"
)

// dice expressions correctly display/eval
func TestDiceExpression(t *testing.T) {
  rand.Seed(time.Now().UTC().UnixNano())

  count,sides := 3,6

  die := NewDiceExpression(count,sides)

  if die.DisplayString() != fmt.Sprintf("%dd%d", count, sides) {
    t.Fail()
  }

  for i := 0; i < 10000; i++ {
    val := die.Eval()
    if val < count || val > count*sides {
      t.Fail()
    }
  }

  if die.Literal() != false {
    t.Fail()
  }
}

// literal expressions correctly display/eval
func TestLiteralExpression(t *testing.T) {
  val := 4

  literal := NewLiteralExpression(val)

  if literal.DisplayString() != fmt.Sprintf("%d",val) ||
  literal.Eval() != val ||
  literal.Literal() != true {
    t.Fail()
  }
}

// addition expressions do what they advertise
func TestAdditionExpression(t *testing.T) {
  lhs, rhs := NewLiteralExpression(3), NewLiteralExpression(6)

  sum := NewAdditionExpression(lhs, rhs)

  if sum.DisplayString() !=
  fmt.Sprintf("%v + %v", lhs.DisplayString(), rhs.DisplayString()) {
    t.Fail()
  }

  if sum.Eval() != lhs.Eval() + rhs.Eval() {
    t.Fail()
  }

  if sum.Literal() != (rhs.Literal() && rhs.Literal()) {
    t.Fail()
  }
}

func checkDiceParse(expr string, disp string, lowBound int, highBound int) bool {
  dExpr, err := ParseDiceExpression(expr)
  if err != nil {
    return false
  }
  if dExpr.DisplayString() != disp {
    return false
  }
  for i := 0; i < 10000; i++ {
    val := dExpr.Eval()
    if val < lowBound || val > highBound {
      return false
    }
  }
  return true
}

func TestDiceExprParse(t *testing.T) {
  if !checkDiceParse("1d4+3", "1d4 + 3", 4, 7) {
    t.Fail()
  }
  if !checkDiceParse("1d4+2d6+3 + \t   6", "1d4 + 2d6 + 9", 12, 25) {
    t.Fail()
  }
  if !checkDiceParse("22+2+3d6", "3d6 + 24", 27, 42) {
    t.Fail()
  }
}

