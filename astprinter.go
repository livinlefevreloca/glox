package main

import (
	"fmt"
	"strings"
)

type AstPrinter struct {
	depth int
	env   Environment
}

func NewAstPrinter(existingEnv *map[string]any) AstPrinter {
	var env map[string]any
	if existingEnv == nil {
		env = make(map[string]any)
	} else {
		env = *existingEnv
	}
	return AstPrinter{depth: 0, env: Environment{values: env}}
}

func (a AstPrinter) print(stmts []Statement) error {
	for _, stmt := range stmts {
		out, err := stmt.accept(a)
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println(out)
	}
	return nil
}

func (a AstPrinter) visitVarDeclarationStatement(stmt VarDeclarationStatement) (any, error) {
	a.depth++
	expr, err := stmt.initializer.accept(a)
	if err != nil {
		return "", err
	}

	out := fmt.Sprintf(
		"VarDeclarationStatement: \n%sName: %s\n%sValue: %s",
		strings.Repeat("\t", a.depth),
		stmt.name,
		strings.Repeat("\t", a.depth),
		expr,
	)

	a.depth--
	return out, nil
}

func (a AstPrinter) visitExpressionStatemet(stmt ExpressionStatement) (any, error) {
	a.depth++
	expr, err := stmt.expr.accept(a)
	if err != nil {
		return "", err
	}

	out := fmt.Sprintf("ExpressionStatement: \n%s%s", strings.Repeat("\t", a.depth), expr)

	a.depth--
	return out, nil
}

func (a AstPrinter) visitPrintStatement(stmt PrintStatement) (any, error) {
	a.depth++
	expr, err := stmt.expr.accept(a)
	if err != nil {
		return "", err
	}

	out := fmt.Sprintf("PrintStatement: \n%s%s", strings.Repeat("\t", a.depth), expr)

	a.depth--
	return out, nil
}

func (a AstPrinter) visitAssign(expr Assign) (any, error) {

	value, err := expr.value.accept(a)
	if err != nil {
		return nil, err
	}

	return fmt.Sprintf("Assign: %s = %s", expr.name.lexeme, value), nil

}

func (a AstPrinter) visitVariable(expr Variable) (any, error) {
	value, err := a.env.get(expr.name.lexeme)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Variable: %s = %v", expr.name.lexeme, value), nil
}

func (a AstPrinter) visitTernary(expr Ternary) (any, error) {
	a.depth++

	condition, err := expr.condition.accept(a)
	if err != nil {
		return "", err
	}

	left, err := expr.left.accept(a)
	if err != nil {
		return "", err
	}

	right, err := expr.right.accept(a)
	if err != nil {
		return "", err
	}

	out := fmt.Sprintf("Ternary: %s", condition) + fmt.Sprintf(
		"\n%sLeft  -> %s", strings.Repeat("\t", a.depth), left) + fmt.Sprintf(
		"\n%sRight -> %s", strings.Repeat("\t", a.depth), right)
	a.depth--
	return out, nil
}

func (a AstPrinter) visitBinary(expr Binary) (any, error) {
	a.depth++

	left, err := expr.left.accept(a)
	if err != nil {
		return "", err
	}

	right, err := expr.right.accept(a)
	if err != nil {
		return "", err
	}

	out := fmt.Sprintf("Binary: %s", expr.operator.operator.lexeme) + fmt.Sprintf(
		"\n%sLeft  -> %s", strings.Repeat("\t", a.depth), left) + fmt.Sprintf(
		"\n%sRight -> %s", strings.Repeat("\t", a.depth), right)
	a.depth--
	return out, nil
}

func (a AstPrinter) visitGrouping(expr Grouping) (any, error) {
	a.depth++

	grouping, err := expr.expression.accept(a)
	if err != nil {
		return "", nil
	}

	out := fmt.Sprintf("Grouping: (\n%s%s\n", strings.Repeat("\t", a.depth), grouping)
	a.depth--
	out += fmt.Sprintf("%s)", strings.Repeat("\t", a.depth))
	return out, nil
}

func (a AstPrinter) visitLiteral(expr Literal) (any, error) {
	return fmt.Sprintf("Literal: %s", expr.value.lexeme), nil
}

func (a AstPrinter) visitOperator(expr Operator) (any, error) {
	return fmt.Sprintf("Operator: %s", expr.operator.lexeme), nil
}

func (a AstPrinter) visitUnary(expr Unary) (any, error) {
	unary, err := expr.right.accept(a)
	if err != nil {
		return "", nil
	}

	return fmt.Sprintf("Unary: %s Right -> %s", expr.operator.lexeme, unary), nil
}
