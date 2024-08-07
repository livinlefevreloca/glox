package main

import (
	"fmt"
	"os"
)

func isTruthy(value any) bool {
	// If the value is nil return fals
	if value == nil {
		return false
	}
	// if the value is boolean return it
	if boolean, ok := value.(bool); ok {
		return boolean
	}

	// Otherwise return true
	return true
}

type Interpreter struct {
	errorContext string
}

func (i *Interpreter) interpert(expr Expr) {
	value, err := i.evaluate(expr)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Printf("%v\n", value)
}

func (i *Interpreter) visitLiteral(expr Literal) (any, error) {
	return expr.value.literal, nil
}

func (i *Interpreter) visitGrouping(expr Grouping) (any, error) {
	return i.evaluate(expr.expression)
}

func (i *Interpreter) visitUnary(expr Unary) (any, error) {
	right, err := i.evaluate(expr.right)
	if err != nil {
		return err, nil
	}

	switch expr.operator.tokenType {
	case TOKEN_BANG:
		return !isTruthy(right), nil
	case TOKEN_MINUS:
		if right, ok := right.(float64); ok {
			return -right, nil
		}
	default:
		return fmt.Errorf(
			"[line: %d, col: %d] Unknown unary operator: %s",
			expr.operator.line,
			expr.operator.col,
			expr.operator.lexeme,
		), nil
	}

	return fmt.Errorf(
		"[line: %d, col: %d] Unexpected values for operator: %s",
		expr.operator.line,
		expr.operator.col,
		expr.operator.lexeme,
	), nil
}

func (i *Interpreter) visitBinary(expr Binary) (any, error) {
	left, err := i.evaluate(expr.left)
	if err != nil {
		return nil, err
	}
	right, err := i.evaluate(expr.right)
	if err != nil {
		return nil, err
	}

	switch expr.operator.operator.tokenType {
	case TOKEN_BANG_EQUAL:
		return left != right, nil
	case TOKEN_EQUAL_EQUAL:
		return left == right, nil
	case TOKEN_GREATER:
		left, okLeft := left.(float64)
		right, okRight := right.(float64)
		if okLeft && okRight {
			return left > right, nil
		}
	case TOKEN_GREATER_EQUAL:
		left, okLeft := left.(float64)
		right, okRight := right.(float64)
		if okLeft && okRight {
			return left >= right, nil
		}
	case TOKEN_LESS:
		left, okLeft := left.(float64)
		right, okRight := right.(float64)
		if okLeft && okRight {
			return left < right, nil
		}
	case TOKEN_LESS_EQUAL:
		left, okLeft := left.(float64)
		right, okRight := right.(float64)
		if okLeft && okRight {
			return left <= right, nil
		}
	case TOKEN_MINUS:
		left, okLeft := left.(float64)
		right, okRight := right.(float64)
		if okLeft && okRight {
			return left - right, nil
		}
	case TOKEN_SLASH:
		left, okLeft := left.(float64)
		right, okRight := right.(float64)
		if okLeft && okRight {
			return left / right, nil
		}
	case TOKEN_STAR:
		left, okLeft := left.(float64)
		right, okRight := right.(float64)
		if okLeft && okRight {
			return left / right, nil
		}
	case TOKEN_PLUS:
		leftStr, okLeftStr := left.(string)
		rightStr, okRightStr := right.(string)
		leftNum, okNumLeftNum := left.(float64)
		rightNum, okNumRightNum := right.(float64)

		if okLeftStr && okRightStr {
			return leftStr + rightStr, nil
		} else if okNumLeftNum && okNumRightNum {
			return leftNum + rightNum, nil
		}
	default:
		return fmt.Errorf(
			"[line: %d, col: %d] Unknown operator: %s",
			expr.operator.operator.line,
			expr.operator.operator.col,
			expr.operator.operator.lexeme,
		), nil
	}

	return fmt.Errorf(
		"[line: %d, col: %d] Unexpected values for operator: %s",
		expr.operator.operator.line,
		expr.operator.operator.col,
		expr.operator.operator.lexeme,
	), nil
}

func (i *Interpreter) visitOperator(expr Operator) (any, error) {
	return nil, nil
}

func (i *Interpreter) visitTernary(expr Ternary) (any, error) {
	condition, err := i.evaluate(expr.condition)
	if err != nil {
		return nil, err
	}

	if isTruthy(condition) {
		return i.evaluate(expr.left)
	} else {
		return i.evaluate(expr.right)
	}
}

func (i *Interpreter) evaluate(expr Expr) (any, error) {
	return expr.accept(i)
}
