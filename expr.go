package main

import (
	"fmt"
	"strings"
)

type AstPrinter struct {
	depth int
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

type Visitor interface {
	visitTernary(expr Ternary) (any, error)
	visitBinary(expr Binary) (any, error)
	visitGrouping(expr Grouping) (any, error)
	visitLiteral(expr Literal) (any, error)
	visitOperator(expr Operator) (any, error)
	visitUnary(expr Unary) (any, error)
}

type Expr interface {
	accept(visitor Visitor) (any, error)
}

type Ternary struct {
	condition Expr
	left      Expr
	right     Expr
}

func (t Ternary) accept(visitor Visitor) (any, error) {
	return visitor.visitTernary(t)
}

type Binary struct {
	left     Expr
	operator Operator
	right    Expr
}

func (b Binary) accept(visitor Visitor) (any, error) {
	return visitor.visitBinary(b)
}

type Grouping struct {
	expression Expr
}

func (g Grouping) accept(visitor Visitor) (any, error) {
	return visitor.visitGrouping(g)
}

type Literal struct {
	value Token
}

func (l Literal) accept(visitor Visitor) (any, error) {
	return visitor.visitLiteral(l)
}

type Operator struct {
	operator Token
}

func (o Operator) accept(visitor Visitor) (any, error) {
	return visitor.visitOperator(o)
}

type Unary struct {
	operator Token
	right    Expr
}

func (u Unary) accept(visitor Visitor) (any, error) {
	return visitor.visitUnary(u)
}
