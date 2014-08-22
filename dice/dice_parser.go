package dice

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

func consumeValue(data []byte, atEOF bool) (int, []byte, error) {
	var accum []byte
	if atEOF {
		return len(data), data, nil
	}
	for i, b := range data {
		if b == ' ' || b == '\n' || b == '\t' || b == '\r' || b == '+' || b == '*' {
			return i, accum, nil
		} else {
			accum = append(accum, b)
		}
	}
	return 0, nil, nil
}

func consumeWhitespace(data []byte, atEOF bool) (int, []byte, error) {
	var accum []byte
	if atEOF {
		return len(data), data, nil
	}
	for i, b := range data {
		if b == ' ' || b == '\n' || b == '\t' || b == '\r' {
			accum = append(accum, b)
		} else {
			return i, accum, nil
		}
	}
	return 0, nil, nil
}

func ParseDiceExpression(expr string) (DiceExpression, error) {
	s := bufio.NewScanner(strings.NewReader(expr))

	split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		switch data[0] {
		case '+', '*':
			advance, token, err = 1, data[:1], nil
		case ' ', '\n', '\r', '\t':
			advance, token, err = consumeWhitespace(data, atEOF)
		default:
			advance, token, err = consumeValue(data, atEOF)
		}
		return
	}
	s.Split(split)

	sumVals := make([]DiceExpression, 0)

	for s.Scan() {
		token := string(s.Bytes())

		// ignore operators (for now) and whitespace
		if token == "+" {
			continue
		} else if token[0] == ' ' || token[0] == '\n' || token[0] == '\t' || token[0] == '\r' {
			continue
		}

		diceSplit := strings.Split(token, "d")

		// it's a literal or operator
		if len(diceSplit) == 1 {
			num, err := strconv.Atoi(diceSplit[0])
			if err != nil {
				return nil, err
			}
			sumVals = append(sumVals, NewLiteralExpression(num))
			// it's a diceExpr
		} else if len(diceSplit) == 2 {
			count, err := strconv.Atoi(diceSplit[0])
			if err != nil {
				return nil, err
			}
			sides, err := strconv.Atoi(diceSplit[1])
			if err != nil {
				return nil, err
			}
			sumVals = append(sumVals, NewDiceExpression(count, sides))
			// it's an error
		} else {
			// TODO: error
			return nil, fmt.Errorf("Can't parse %v", diceSplit)
		}
	}

	sumLiteral := NewLiteralExpression(0)
	sumRemaining := (DiceExpression)(nil)
	for _, dExpr := range sumVals {
		if dExpr.Literal() {
			sumLiteral = NewLiteralExpression(sumLiteral.Eval() + dExpr.Eval())
		} else {
			if sumRemaining == nil {
				sumRemaining = dExpr
			} else {
				sumRemaining = NewAdditionExpression(sumRemaining, dExpr)
			}
		}
	}
	if sumRemaining == nil {
		return sumLiteral, nil
	} else {
		return NewAdditionExpression(sumRemaining, sumLiteral), nil
	}
}
