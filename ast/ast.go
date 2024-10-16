package ast

import (
	"bytes"
	"github.com/ldcicconi/monkey-interpreter/token"
	"strings"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

type LetStatement struct {
	Token token.Token // the token.LET token
	Name  Identifier
	Value Expression
}

func (ls LetStatement) statementNode()       {}
func (ls LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")

	return out.String()

}

type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i Identifier) expressionNode()      {}
func (i Identifier) TokenLiteral() string { return i.Token.Literal }
func (i Identifier) String() string       { return i.Value }

type ReturnStatement struct {
	Token       token.Token // the 'return' token
	ReturnValue Expression
}

func (rs ReturnStatement) statementNode()       {}
func (rs ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}

type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (es ExpressionStatement) statementNode()       {}
func (es ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}

	return ""
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il IntegerLiteral) expressionNode()      {}
func (il IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il IntegerLiteral) String() string       { return il.Token.Literal }

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe PrefixExpression) expressionNode()      {}
func (pe PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe PrefixExpression) String() string {
	var builder strings.Builder

	builder.WriteString("(")
	builder.WriteString(pe.Operator)
	builder.WriteString(pe.Right.String())
	builder.WriteString(")")

	return builder.String()
}

type InfixExpression struct {
	Token       token.Token // operator token, eg "+"
	Left, Right Expression
	Operator    string
}

func (ie InfixExpression) expressionNode()      {}
func (ie InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie InfixExpression) String() string {
	var builder strings.Builder

	builder.WriteString("(")
	builder.WriteString(ie.Left.String())
	builder.WriteString(" " + ie.Operator + " ")
	builder.WriteString(ie.Right.String())
	builder.WriteString(")")

	return builder.String()
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (b Boolean) expressionNode()      {}
func (b Boolean) TokenLiteral() string { return b.Token.Literal }
func (b Boolean) String() string       { return b.Token.Literal }

type IfExpression struct {
	Token       token.Token // "if" token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie IfExpression) expressionNode()      {}
func (ie IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie IfExpression) String() string {
	var builder strings.Builder

	builder.WriteString("if")
	builder.WriteString(ie.Condition.String())
	builder.WriteString(" ")
	builder.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		builder.WriteString("else ")
		builder.WriteString(ie.Alternative.String())
	}

	return builder.String()
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs BlockStatement) statementNode()       {}
func (bs BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs BlockStatement) String() string {
	var builder strings.Builder

	for _, stmt := range bs.Statements {
		builder.WriteString(stmt.String())
	}

	return builder.String()
}

type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl FunctionLiteral) expressionNode()      {}
func (fl FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl FunctionLiteral) String() string {
	var builder strings.Builder
	params := make([]string, 0, len(fl.Parameters))
	for _, param := range fl.Parameters {
		params = append(params, param.String())
	}

	builder.WriteString(fl.TokenLiteral())
	builder.WriteString("(")
	builder.WriteString(strings.Join(params, ", "))
	builder.WriteString(")")
	builder.WriteString(fl.Body.String())

	return builder.String()
}

type CallExpression struct {
	Token     token.Token
	Function  Expression // Identifier or function literal
	Arguments []Expression
}

func (ce CallExpression) expressionNode()      {}
func (ce CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce CallExpression) String() string {
	var builder strings.Builder
	var args []string

	for _, arg := range ce.Arguments {
		args = append(args, arg.String())
	}
	builder.WriteString(ce.Function.String())
	builder.WriteString("(")
	builder.WriteString(strings.Join(args, ", "))
	builder.WriteString(")")

	return builder.String()
}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl StringLiteral) expressionNode()      {}
func (sl StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl StringLiteral) String() string       { return sl.Token.Literal }

type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
}

func (al ArrayLiteral) expressionNode()      {}
func (al ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al ArrayLiteral) String() string {
	var builder strings.Builder

	elements := make([]string, 0, len(al.Elements))
	for _, e := range al.Elements {
		elements = append(elements, e.String())
	}
	joined := strings.Join(elements, ", ")

	builder.WriteString("[")
	builder.WriteString(joined)
	builder.WriteString("]")

	return builder.String()
}

type IndexExpression struct {
	Token token.Token
	Left  Expression
	Index Expression
}

func (i IndexExpression) expressionNode()      {}
func (i IndexExpression) TokenLiteral() string { return i.Token.Literal }
func (i IndexExpression) String() string {
	var builder strings.Builder

	builder.WriteString("(")
	builder.WriteString(i.Left.String())
	builder.WriteString("[")
	builder.WriteString(i.Index.String())
	builder.WriteString("])")

	return builder.String()
}

type HashLiteral struct {
	Token token.Token
	Pairs map[Expression]Expression
}

func (hl HashLiteral) expressionNode()      {}
func (hl HashLiteral) TokenLiteral() string { return hl.Token.Literal }
func (hl HashLiteral) String() string {
	var builder strings.Builder

	pairs := make([]string, 0, len(hl.Pairs))
	for key, value := range hl.Pairs {
		pairs = append(pairs, key.String()+":"+value.String())
	}

	builder.WriteString("{")
	builder.WriteString(strings.Join(pairs, ", "))
	builder.WriteString("}")

	return builder.String()
}
