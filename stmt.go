package main

type StatementVisitor interface {
	visitExpressionStatemet(stmt ExpressionStatement) (any, error)
	visitPrintStatement(stmt PrintStatement) (any, error)
	visitVarDeclarationStatement(stmt VarDeclarationStatement) (any, error)
}

type Statement interface {
	accept(visitor StatementVisitor) (any, error)
}

type VarDeclarationStatement struct {
	name        Token
	initializer Expr
}

func (v VarDeclarationStatement) accept(visitor StatementVisitor) (any, error) {
	return visitor.visitVarDeclarationStatement(v)
}

type ExpressionStatement struct {
	expr Expr
}

func (e ExpressionStatement) accept(visitor StatementVisitor) (any, error) {
	return visitor.visitExpressionStatemet(e)
}

type PrintStatement struct {
	expr Expr
}

func (p PrintStatement) accept(visitor StatementVisitor) (any, error) {
	return visitor.visitPrintStatement(p)
}
