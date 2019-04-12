package lexer

import (
	"github.com/shouji-kazuo/gomonkey/monkey/token"
)

// Lexer is ...
type Lexer struct {
	input        string
	position     int  // 入力の中の現在の位置(現在の文字を指し示す)．現在の文字 == ch．
	readPosition int  // これから読み込む位置(現在の文字の次)
	ch           byte // 現在検査中の文字
}

// New is ...
func New(input string) *Lexer {
	l := &Lexer{
		input: input,
	}
	l.readChar()
	return l
}

// NextToken is ...
// Lexerの「現在の文字」からトークンを生成して返す
// 副作用として，Lexerの「現在の文字」を1文字進める
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			// "==" を処理する
			ch := l.ch
			l.readChar() // 1文字読み飛ばす('='と分かっている->字句解析の状態を進める)
			lit := string(ch) + string(l.ch)
			tok = token.Token{
				Type:    token.EQ,
				Literal: lit,
			}
			// ここは早期リターンできない -> さらに字句解析の状態を進める必要があるため
			// (上のl.readChar()は「すでに分かっている文字」を読み勧めたにすぎず，その先にある文字を1文字読み進める必要がある)
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		if l.peekChar() == '=' {
			// "!=" を処理する
			ch := l.ch
			l.readChar() // 1文字読み飛ばす('='と分かっている->字句解析の状態を進める)
			lit := string(ch) + string(l.ch)
			tok = token.Token{
				Type:    token.NOT_EQ,
				Literal: lit,
			}
			// ここは早期リターンできない -> さらに字句解析の状態を進める必要があるため
			// (上のl.readChar()は「すでに分かっている文字」を読み勧めたにすぎず，その先にある文字を1文字読み進める必要がある)
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
			// なぜreadChar()により「現在の文字」を更新せずreturnするか
			// readIdentifier()の中ですでに更新されていることから，
			// このあとの l.readChar()で余分に更新されることを防ぐため
		}
		if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		}
		tok = newToken(token.ILLEGAL, l.ch)
	}

	l.readChar()
	return tok
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) skipWhitespace() {
	ch := l.ch
	for ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
		ch = l.ch
	}
}

func (l *Lexer) readIdentifier() string {
	position := l.position // 今の「現在の文字の位置」を保存
	for isLetter(l.ch) {
		l.readChar()
	}

	// 保存しておいた「現在の文字の位置」から，「isLetter()がtrueとなる最後の位置」までの
	// 部分文字列を識別子として返す
	return l.input[position:l.position]
}

// isLetter is ...
// 与えられた文字が英字('_'を含)かどうか
// (識別子として使える文字かどうか))
func isLetter(ch byte) bool {
	return ('a' <= ch && ch <= 'z') ||
		('A' <= ch && ch <= 'Z') ||
		ch == '_'
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string(ch),
	}
}

// 入力から1文字読む．字句解析の状態を更新しない．
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

// 入力から1文字読む．字句解析の状態を更新する．
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		// これから読み込もうとする文字の位置が入力の長さを超えていたら
		// NULを入れる
		// (NUL == 「まだ何も読み込んでいない」 or 「ファイルの終わり」) ...? EOFじゃなくて?
		l.ch = 0
	} else {
		// これから読み込もうとする文字の位置が入力の長さ内なら
		// 文字を読む
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}
