package switchparser

import (
	"slices"
	"strings"

	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/lexer"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/token"
)

func ToString(node ast.Node) string {
	var b strings.Builder
	node.WriteTo(&b)
	return b.String()
}

type SwitchStatement struct {
	discriminant ast.Expression
	cases        []SwitchCase
}

type SwitchCase struct {
	test       ast.Expression
	consequent []ast.Statement
}

func (ss *SwitchStatement) WriteTo(b *strings.Builder) {
	b.WriteString("switch(")
	ss.discriminant.WriteTo(b)
	b.WriteString("){")
	for i, stmt := range ss.cases {
		if i > 0 {
			b.WriteRune(';')
		}
		stmt.WriteTo(b)
	}
	b.WriteRune('}')
}

func (sc *SwitchCase) WriteTo(b *strings.Builder) {
	if sc.test != nil {
		b.WriteString("case ")
		sc.test.WriteTo(b)
		b.WriteRune(':')
	} else {
		b.WriteString("default:")
	}
	for i, stmt := range sc.consequent {
		if i > 0 {
			b.WriteRune(';')
		}
		stmt.WriteTo(b)
	}
}

func Plugin(pb *parser.Builder) {
	lb := pb.LexerBuilder

	// registers keywords
	switchTokenType := lb.RegisterTokenType("switch")
	switchCaseTokenType := lb.RegisterTokenType("case")
	defaultTokenType := lb.RegisterTokenType("default")
	lb.UseTokenInterceptor(func(l *lexer.Lexer, next func() token.Token) token.Token {
		ret := next()
		if ret.Type != token.IDENT {
			return ret
		}
		switch ret.Literal {
		case "switch":
			ret.Type = switchTokenType
		case "case":
			ret.Type = switchCaseTokenType
		case "default":
			ret.Type = defaultTokenType
		}
		return ret
	})

	// intercepts the statement
	pb.UseStatementInterceptor(func(p *parser.Parser, next func() ast.Statement) ast.Statement {
		if p.CurrentToken.Type != switchTokenType {
			return next()
		}
		switchStmt := &SwitchStatement{
			cases: []SwitchCase{},
		}
		if !p.ExpectToken(token.LPAREN) {
			return nil
		}
		p.NextToken() // move to expression
		switchStmt.discriminant = p.ParseExpression()
		if !p.ExpectToken(token.RPAREN) {
			return nil
		}
		if !p.ExpectToken(token.LBRACE) {
			return nil
		}
		p.NextToken()
		caseTypes := []token.Type{switchCaseTokenType, defaultTokenType}
		for slices.Contains(caseTypes, p.CurrentToken.Type) {
			switchCase := SwitchCase{}
			if p.CurrentToken.Type == switchCaseTokenType {
				p.NextToken()
				switchCase.test = p.ParseExpression()
			}
			if !p.ExpectToken(token.COLON) {
				return nil
			}
			p.NextToken()
			for p.CurrentToken.Type != token.RBRACE && !slices.Contains(caseTypes, p.CurrentToken.Type) {
				stmt := p.ParseStatement()
				switchCase.consequent = append(switchCase.consequent, stmt)
				p.NextToken()
			}
			switchStmt.cases = append(switchStmt.cases, switchCase)
		}
		return switchStmt
	})
}
