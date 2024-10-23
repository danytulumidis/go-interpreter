package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
	"strconv"
)

// iota numbers all const variables starting from 1
// the _ consumes the 0
// We need this for precedence
const (
	_ int = iota
	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

	// Hash map to check if a token has a associated parsing function
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

// Define types for the Expression parsing
// Very nice so we can define multiple functions for different tokens
// and store them into our Hash map
type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

// Helper functions to register the right function in the Hash Map
// Second parameter is our type
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func New(l *lexer.Lexer) *Parser {
	// Init the lexer in our Parser with the parameter lexer (pointer so the address of the Lexer object)
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// Use make to initialize a Hash Table to register different expression parsing functions
	// for each token type
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)

	// Read two tokens so curToken AND peekToken are set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	// Advance token from the prefix (e.g. -5, from the - to 5)
	p.nextToken()

	// Parse the e.g. 5
	expression.Right = p.parseExpression(PREFIX)

	return expression
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
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
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

// Parse return statements
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	// Advance token to be on the expression after the =
	p.nextToken()

	// TODO: Were skipping the expressions until we encounter a semicolon
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// Parse expression statement
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	// Get the parsing function from our Hash Table for this token
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	// Parse the left side of our Expression
	leftExp := prefix()

	return leftExp
}

// Returns a AST Identifier with the token and its Value
// DOESNT advance the token
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = value

	return lit
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

// Append an error to our Parser slice
func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}
