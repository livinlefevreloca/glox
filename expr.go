package main

import (
	"fmt"
	"strings"
)

type AstPrinter struct {
	depth int
}

func (a AstPrinter) visitTernary(expr Ternary[string]) string {
	a.depth++
	out := fmt.Sprintf("Ternary: %s", expr.condition.accept(a)) + fmt.Sprintf(
		"\n%sLeft -> %s", strings.Repeat("\t", a.depth), expr.left.accept(a)) + fmt.Sprintf(
		"\n%sRight -> %s", strings.Repeat("\t", a.depth), expr.right.accept(a))
	a.depth--
	return out
}

func (a AstPrinter) visitBinary(expr Binary[string]) string {
	a.depth++
	out := fmt.Sprintf("Binary: %s", expr.operator.operator.lexeme) + fmt.Sprintf(
		"\n%sLeft %s", strings.Repeat("\t", a.depth), expr.left.accept(a)) + fmt.Sprintf(
		"\n%sRight -> %s", strings.Repeat("\t", a.depth), expr.right.accept(a))
	a.depth--
	return out
}

func (a AstPrinter) visitGrouping(expr Grouping[string]) string {
	a.depth++
	out := fmt.Sprintf("Grouping: (\n%s%s\n", strings.Repeat("\t", a.depth), expr.expression.accept(a))
	a.depth--
	out += fmt.Sprintf("%s)", strings.Repeat("\t", a.depth))
	return out
}

func (a AstPrinter) visitLiteral(expr Literal[string]) string {
	return fmt.Sprintf("Literal: %s", expr.value.lexeme)
}

func (a AstPrinter) visitOperator(expr Operator[string]) string {
	return fmt.Sprintf("Operator: %s", expr.operator.lexeme)
}

func (a AstPrinter) visitUnary(expr Unary[string]) string {
	return fmt.Sprintf("Unary: Right -> %s", expr.right.accept(a))
}

type Visitor[T comparable] interface {
	visitTernary(expr Ternary[T]) T
	visitBinary(expr Binary[T]) T
	visitGrouping(expr Grouping[T]) T
	visitLiteral(expr Literal[T]) T
	visitOperator(expr Operator[T]) T
	visitUnary(expr Unary[T]) T
}

type Expr[T comparable] interface {
	accept(visitor Visitor[T]) string
}

type Ternary[T comparable] struct {
	condition Expr[T]
	left      Expr[T]
	right     Expr[T]
}

func (t Ternary[T]) accept(visitor Visitor[T]) T {
	return visitor.visitTernary(t)
}

type Binary[T comparable] struct {
	left     Expr[T]
	operator Operator[T]
	right    Expr[T]
}

func (b Binary[T]) accept(visitor Visitor[T]) T {
	return visitor.visitBinary(b)
}

type Grouping[T comparable] struct {
	expression Expr[T]
}

func (g Grouping[T]) accept(visitor Visitor[T]) T {
	return visitor.visitGrouping(g)
}

type Literal[T comparable] struct {
	value Token
}

func (l Literal[T]) accept(visitor Visitor[T]) T {
	return visitor.visitLiteral(l)
}

type Operator[T comparable] struct {
	operator Token
}

func (o Operator[T]) accept(visitor Visitor[T]) T {
	return visitor.visitOperator(o)
}

type Unary[T comparable] struct {
	operator Token
	right    Expr[T]
}

func (u Unary[T]) accept(visitor Visitor[T]) T {
	return visitor.visitUnary(u)
}
