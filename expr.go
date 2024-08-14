package main

type ExprVisitor interface {
	visitAssign(expr Assign) (any, error)
	visitVariable(expr Variable) (any, error)
	visitTernary(expr Ternary) (any, error)
	visitBinary(expr Binary) (any, error)
	visitGrouping(expr Grouping) (any, error)
	visitLiteral(expr Literal) (any, error)
	visitOperator(expr Operator) (any, error)
	visitUnary(expr Unary) (any, error)
}

type Expr interface {
	accept(visitor ExprVisitor) (any, error)
}

type Assign struct {
	name  Token
	value Expr
}

func (a Assign) accept(visitor ExprVisitor) (any, error) {
	return visitor.visitAssign(a)
}

type Variable struct {
	name Token
}

func (v Variable) accept(visitor ExprVisitor) (any, error) {
	return visitor.visitVariable(v)
}

type Ternary struct {
	condition Expr
	left      Expr
	right     Expr
}

func (t Ternary) accept(visitor ExprVisitor) (any, error) {
	return visitor.visitTernary(t)
}

type Binary struct {
	left     Expr
	operator Operator
	right    Expr
}

func (b Binary) accept(visitor ExprVisitor) (any, error) {
	return visitor.visitBinary(b)
}

type Grouping struct {
	expression Expr
}

func (g Grouping) accept(visitor ExprVisitor) (any, error) {
	return visitor.visitGrouping(g)
}

type Literal struct {
	value Token
}

func (l Literal) accept(visitor ExprVisitor) (any, error) {
	return visitor.visitLiteral(l)
}

type Operator struct {
	operator Token
}

func (o Operator) accept(visitor ExprVisitor) (any, error) {
	return visitor.visitOperator(o)
}

type Unary struct {
	operator Token
	right    Expr
}

func (u Unary) accept(visitor ExprVisitor) (any, error) {
	return visitor.visitUnary(u)
}
