package main

import (
	"fmt"
)

func isTruthy(value any) bool {
	// If the value is nil return false
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
	environment Environment
	scopeDepth  int
}

func NewInterpreter(existingEnv *map[string]any) *Interpreter {

	var env map[string]any
	if existingEnv == nil {
		env = make(map[string]any)
	} else {
		env = *existingEnv
	}

	return &Interpreter{
		environment: Environment{name: "INTENV_BASE", values: env},
	}
}

func (i *Interpreter) interpert(statements []Statement) (any, error) {
	var result any
	var err error
	for _, statement := range statements {
		result, err = i.execute(statement)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (i *Interpreter) execute(stmt Statement) (any, error) {
	return stmt.accept(i)
}

func (i *Interpreter) visitExpressionStatemet(stmt ExpressionStatement) (any, error) {
	return i.evaluate(stmt.expr)
}

func (i *Interpreter) visitVarDeclarationStatement(stmt VarDeclarationStatement) (any, error) {
	var value any
	var err error
	if stmt.initializer != nil {
		value, err = i.evaluate(stmt.initializer)
		if err != nil {
			return nil, err
		}
	}

	i.environment.define(stmt.name.lexeme, value)
	return nil, nil
}

func (i *Interpreter) visitBlockStatement(stmt BlockStatement) (any, error) {
	i.scopeDepth++
	i.executeBlock(stmt.stmts, Environment{name: fmt.Sprintf("INTENV_%d", i.scopeDepth), values: make(map[string]any)})
	i.scopeDepth--
	return nil, nil
}

func (i *Interpreter) executeBlock(stmts []Statement, env Environment) error {
	previousEnv := i.environment

	env.parent = &previousEnv
	i.environment = env

	for _, stmt := range stmts {
		_, err := i.execute(stmt)
		if err != nil {
			return err
		}
	}

	i.environment = previousEnv

	return nil
}

func (i *Interpreter) visitPrintStatement(stmt PrintStatement) (any, error) {
	value, err := i.evaluate(stmt.expr)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%v\n", value)
	return nil, nil
}

func (i *Interpreter) visitAssign(expr Assign) (any, error) {
	value, err := i.evaluate(expr.value)
	if err != nil {
		return nil, err
	}
	err = i.environment.assign(expr.name.lexeme, value)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (i *Interpreter) visitVariable(expr Variable) (any, error) {
	value, err := i.environment.get(expr.name.lexeme)
	if err != nil {
		return nil, err
	}

	return value, nil
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
	case TOKEN_COMMA:
		return right, nil
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

		if right == 0 {
			return fmt.Errorf(
				"[line: %d, col: %d] Division by zero",
				expr.operator.operator.line,
				expr.operator.operator.col,
			), nil
		}

		if okLeft && okRight {
			return left / right, nil
		}
	case TOKEN_STAR:
		left, okLeft := left.(float64)
		right, okRight := right.(float64)

		if okLeft && okRight {
			return left * right, nil
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
