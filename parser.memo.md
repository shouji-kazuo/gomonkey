# WIP
# 概要
## パースの流れ(`parseExpression`)

1. 単項式と思ってパースし始める
    * <単項演算子>?<リテラル>
2. `peekToken`に2項演算子がきた && `parseExpression`に与えられた優先度より`peekToken`の優先度が高ければ，
    1. トークンを1つ進めて，
    2. `peekToken`以降を多項式の第2項以降とみなしてパースを試みる
        1. 第2項以降に，今の `peekToken` より優先順位が高い演算子が来た場合
        2. 


# トレース

## 例: 1 * 2 + 3;

```
parser.New()
    nextToken()
        curToken = null
        peekToken = 1
    nextToken()
        curToken = 1
        peekToken = +
ParseProgram()
^^^^^^^^^^^^
    statement = parseStatement()
                ^^^^^^^^^^^^^^^
        return parseExpressionStatement()
               ^^^^^^^^^^^^^^^^^^^^^^^^^
            statement = ExpressionStatement {
                Token: 1
                Expression: parseExpression(LOWEST)
                            ^^^^^^^^^^^^^^^^^^^^^^^
            }
                prefix = parseIntegerLiteral
                leftExp = parseIntegerLiteral()
                          ^^^^^^^^^^^^^^^^^^^^^
                    return IntegerLiteral {
                        Token: curToken(1)
                        Value: 1
                    }
                (leftExp に IntegerLiteralを得る)
                for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence()
                    → peekToken = + なので第1項はtrue
                      precedence は LOWEST , peekPrecedenceは SUM なので第2項はtrue
                    → for の中身を実行
                    infix = parseInfixExpression
                    nextToken()
                        curToken = *
                        peekToken = 2
                    leftExp = parseInfixExpression()
                              ^^^^^^^^^^^^^^^^^^^^
                        expression = InfixExpression {
                            Token: *
                            Operator: *
                            Left: IntegerLiteral(1)
                        }
                        precedence = PRODUCT(5)
                        nextToken()
                            curToken: 2
                            peekToken: +
                        expression.Right = parseExpression(PRODUCT(5))
                                           ^^^^^^^^^^^^^^^
                            prefix = parseIntergerLiteral
                            leftExp = parseInteerLiteral()
                                      ^^^^^^^^^^^^^^^^^^
                                return IntegerLiteral {
                                    Token: 2
                                    Value: 2
                                }
                            for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence()
                                → 第1項はtrue
                                  第2項はfalse
                                    precedence = PRODUT(5)
                                    peekPrecedence = SUM(4)
                                → for の中身を実行しない
                            return leftExp
                                leftExp = return IntegerLiteral {
                                    Token: 2
                                    Value: 2
                                }
                        expressionに以下を得る
                            expression = InfixExpression {
                                Token: *
                                Operator: *
                                Left: IntegerLiteral(1)
                                Right: IntegerLiteral {
                                    Token: 2
                                    Value: 2
                                }
                            }
                        ...
```


## 例: 1 + 2 + 3;

```
parser.New()
    nextToken()
        curToken = null
        peekToken = 1
    nextToken()
        curToken = 1
        peekToken = +
ParseProgram()
^^^^^^^^^^^^
    statement = parseStatement()
                ^^^^^^^^^^^^^^^
        return parseExpressionStatement()
               ^^^^^^^^^^^^^^^^^^^^^^^^^
            statement = ExpressionStatement {
                Token: 1
                Expression: parseExpression(LOWEST)
                            ^^^^^^^^^^^^^^^^^^^^^^^
            }
                prefix = parseIntegerLiteral
                leftExp = parseIntegerLiteral()
                          ^^^^^^^^^^^^^^^^^^^^^
                    return IntegerLiteral {
                        Token: curToken(1)
                        Value: 1
                    }
                (leftExp に IntegerLiteralを得る)
                for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence()
                    → peekToken = + なので第1項はtrue
                      precedence は LOWEST , peekPrecedenceは SUM なので第2項はtrue
                    → for の中身を実行
                    infix = parseInfixExpression
                    nextToken()
                        curToken = +
                        peekToken = 2
                    leftExp = parseInfixExpression(leftExp)
                              ^^^^^^^^^^^^^^^^^^^^
                        expression = InfixExpression {
                            Token: +
                            Operator: +
                            Left: IntegerLiteral(1)
                        }
                        precedence = SUM(4)
                        nextToken()
                            curToken = 2
                            peekToken = +
                        expression.Right = parseExpression(SUM)
                                           ^^^^^^^^^^^^^^^
                            prefix = parseIntegerLiteral
                            leftExp = parseIntegerLiteral()
                                      ^^^^^^^^^^^^^^^^^^^^
                                literal := ast.IntegerLiteral {
                                    Token: 2
                                    Value: 2
                                }
                                return literal
                        return expression
                            expressionには以下を得ている
                            InfixExpression {
                                Token: +
                                Operator: +
                                Left: IntegerLiteral(1)
                                Right: IntegerLiteral(2)
                            }
                    leftExpには以下を得ている
                        InfixExpression {
                            Token: +
                            Operator: +
                            Left: IntegerLiteral(1)
                            Right: IntegerLiteral(2)
                        }
                    for の実行条件
                        → 第1項はtrue (peekToken == +)
                          第2項はtrue
                            precedence == LOWEST
                            peekPrecedence == SUM
                        → forの中身を実行する
                    infix = parseInfixExpression (2個目の+により得る)
                    nextToken()
                        curToken = +
                        peekToken = 3
                    leftExp = parseInfixExpression(leftExp)
                                ^^^^^^^^^^^^^^^^^^^^
                        expression = ast.InfixExpression {
                            Token: +
                            Operator: +
                            Left: InfixExpression {
                                Token: +
                                Operator: +
                                Left: IntegerLiteral(1)
                                Right: IntegerLiteral(2)
                            }
                        }
                        nextToken()
                            curToken = 3
                            peekToken = ;
                        expression.Right = parseExpression(PREFIX)
                                           ^^^^^^^^^^^^^^^
                            prefix = parseIntegerLiteral
                            leftExp = parseIntegerLiteral()
                                      ^^^^^^^^^^^^^^^^^^^
                                return ast.IntegerLiteral {
                                    Token: 3
                                    Value: 3
                                }
                            for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence()
                                → 第1項はtrue
                                  第2項はfalse
                                    precedence == PREFIX
                                    peekPrecedence() == LOWEST
                                        ; に対応するprecedenceがないためLOWESTが返る
                                → for は実行されない
                            return leftExp
                        return expression
                            expresssionに以下を得ている
                                expression = ast.InfixExpression {
                                    Token: +
                                    Operator: +
                                    Left: InfixExpression {
                                        Token: +
                                        Operator: +
                                        Left: IntegerLiteral(1)
                                        Right: IntegerLiteral(2)
                                    }
                                    Right: ast.IntegerLiteral {
                                        Token: 3
                                        Value: 3
                                    }
                                }
                    forの実行条件はfalse
                return lextExp
            statement = ExpressionStatement {
                Token: curToken(1)
                Expression: ast.InfixExpression {
                    Token: +
                    Operator: +
                    Left: InfixExpression {
                        Token: +
                        Operator: +
                        Left: IntegerLiteral(1)
                        Right: IntegerLiteral(2)
                    }
                    Right: ast.IntegerLiteral {
                        Token: 3
                        Value: 3
                    }
                }
            }
            return statement
```

## 例1: 1 + 2 * 3;

```
parser.New()
    nextToken()
        curToken = null
        peekToken = 1
    nextToken()
        curToken = 1
        peekToken = +
ParseProgram()
^^^^^^^^^^^^
    statement = parseStatement()
                ^^^^^^^^^^^^^^^
        return parseExpressionStatement()
               ^^^^^^^^^^^^^^^^^^^^^^^^^
            statement = ExpressionStatement {
                Token: curToken(1)
                Expression: parseExpression(LOWEST)
                            ^^^^^^^^^^^^^^^^^^^^^^^
            }
                prefix = parseIntegerLiteral
                leftExp = parseIntegerLiteral()
                          ^^^^^^^^^^^^^^^^^^^^^
                    return IntegerLiteral {
                        Token: curToken(1)
                        Value: 1
                    }
                (leftExp に IntegerLiteralを得る)
                for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence()
                    → peekToken = + なので第1項はtrue
                      precedence は LOWEST , peekPrecedenceは SUM なので第2項はtrue
                    → for の中身を実行
                    infix = parseInfixExpression
                    nextToken()
                        curToken = +
                        peekToken = 2
                    leftExp = parseInfixExpression(leftExp)
                              ^^^^^^^^^^^^^^^^^^^^
                        expression = InfixExpression {
                            Token: +
                            Operator: +
                            Left: IntegerLiteral(1)
                        }
                        precedence = SUM
                        nextToken()
                            curToken = 2
                            peekToken = *
                        expression.Right = parseExpression(SUM)
                                           ^^^^^^^^^^^^^^^
                            prefix = parseIntegerLiteral
                            leftExp = parseIntegerLiteral()
                                      ^^^^^^^^^^^^^^^^^^^
                                return IntegerLiteral {
                                    Token: 2
                                    Value: 2
                                }
                            for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence()
                                → peekToken = * なので第1項はtrue
                                  precedence = SUM(4), peekPrecedence = ASTERLISK(4) なので 第2項はtrue
        
                                → forの中身を実行
                                infix = parseInfixExpression
                                nextToken()
                                    curToken = *
                                    peekToken = 3
                                leftExp = parseInfixExpression(IntegerLiteral(2))
                                          ^^^^^^^^^^^^^^^^^^^^
                                    expression = InfixExpression {
                                        Token: *
                                        Operator: *
                                        Left: IntegerLiteral(2)
                                    }
                                    precedence = ASTERLISK(4)
                                    nextToken()
                                        curToken = 3
                                        peekToken = ;
                                    expression.Right = parseExpression(ASTERLISK(4))
                                                       ^^^^^^^^^^^^^^^
                                        prefix = parseIntegerLiteral
                                        leftExp = parseIntegerLiteral()
                                            return IntegerLiteral {
                                                Token: 3
                                                Value: 3
                                            }
                                        for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence()
                                            → peekToken = ; なので第1項はfalse
                                            → forの中身は実行されない
                                        return leftExp ... IntegerLiteral(3)
                                    expression に以下を得る
                                        expression = {
                                            Token: *
                                            Operator: *
                                            Left: IntegerLiteral(2)
                                            Right: IntegerLiteral(3)
                                        }
                                        2 * 3 を得る
                                        ```
                                           *
                                          / \
                                         2   3
                                        ```
                                    return expression
                                leftExp = 上記のexpression
                                peekToken = ; なのでforの実行条件を満たさない
                                forを抜ける
                                return leftExp(2 * 3)
                        expression.Right に 2 * 3 の式を得る
                        return expression {
                            Token: +
                            Operator: +
                            Left: IntegerLiteral(1)
                            Right: 2 * 3
                        }
                    leftExp = 上記のexpression
                    peekToken = ; なのでforの実行条件を満たさない
                    forを抜ける
                return leftExp
            statements = 上記のleftExp
            nextToken()
                curToken: ;
                peekToken = EOF
            終わり
```