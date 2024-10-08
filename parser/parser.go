package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
)

type Parser struct {
	l *lexer.Lexer

	errors []string

	curToken  token.Token
	peekToken token.Token
}

func New(l *lexer.Lexer) *Parser {
	// Init the lexer in our Parser with the parameter lexer (pointer so the address of the Lexer object)
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// Read two tokens so curToken AND peekToken are set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

// Helper function to set the current and next token (similar to position and readPosition in our Lexer)
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// Entry point of our parser, init our AST and set the statements to an empty slice
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.curTokenIs(token.EOF) {
		// Parse each statement
		stmt := p.parseStatement()
		if stmt != nil {
			// Append the statement to our AST program statements slice
			program.Statements = append(program.Statements, stmt)
		}
		// Move to the next Token
		p.nextToken()
	}

	return program
}

// Parse each statement and return it
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	default:
		return nil
	}
}

// Parse a Let Statement (e.g. let x = 5)
func (p *Parser) parseLetStatement() *ast.LetStatement {
	// Creates a new Statement pointer to our LetStatement struct from the AST
	// Init the Token field (LET token)
	stmt := &ast.LetStatement{Token: p.curToken}

	// Check if next statement is an identifier (e.g. x or foo), if not its not a valid let statement
	// expectPeek moves to the next token
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	// The name which is an Identifier now points to a AST Identifier struct with all fields initialized
	// Literal would be x or foo etc.
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// Next Token should be an assign token
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO: Skipping the expression until we encounter a semicolon
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// Helper function to validate a token for a specific type
func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

// Helper function to validate the next token with a specific type
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// Check a token for a specific type and then advance to the next token
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}
