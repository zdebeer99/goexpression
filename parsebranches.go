package goexpression

func (this *parser) pumpExpression() {
	this.state = branchExpressionValuePart
	for this.state != nil {
		if this.err != nil {
			break
		}
		this.state = this.state(this)
	}
	endo := this.commit()
	if len(endo) > 0 || !this.scan.IsEOF() {
		this.error("Unexpected end of expression. '" + endo + "' not parsed. ")
	}
}

/*
parse expressions
[value part][operator part] repeat
*/

//
func branchExpressionValuePart(this *parser) stateFn {
	scan := this.scan
	scan.SkipSpaces()
	if scan.IsEOF() {
		return nil
	}
	if scan.ScanNumber() {
		this.add(NewNumberToken(scan.Commit()))
		return branchExpressionOperatorPart
	}
	if scan.ScanWord() {
		return branchExpressionAfterWord
	}
	switch scan.Next() {
	case '"', '\'':
		scan.Backup()
		txt := this.ParseText()
		this.add(NewTextToken(txt))
		return branchExpressionOperatorPart
	case '(':
		this.parseOpenBracket()
		return branchExpressionValuePart
	}
	this.error("Unexpected token. ")
	return nil
}

func branchExpressionAfterWord(this *parser) stateFn {
	scan := this.scan
	switch scan.Peek() {
	case '(':
		this.curr = this.add(NewFuncToken(scan.Commit()))
		return branchFunctionArguments
	}
	this.add(NewIdentityToken(scan.Commit()))
	return branchExpressionOperatorPart
}

func branchFunctionArguments(this *parser) stateFn {
	scan := this.scan
	r := scan.Next()
	if r != '(' {
		this.error("Expecting '(' before arguments.")
	}
	ftoken, ok := this.curr.Value.(*FuncToken)
	if !ok {
		this.error("Expecting function token to add arguments to.")
		return nil
	}
	state := branchExpressionValuePart
	currnode := this.curr
	for {
		this.curr = NewTreeNode(NewGroupToken(""))
		for state != nil {
			state = state(this)
		}
		r = scan.Next()
		switch r {
		case ' ':
			scan.Ignore()
			continue
		case ',':
			ftoken.AddArgument(this.curr.Root())
			state = branchExpressionValuePart
			scan.Ignore()
			continue
		case ')':
			ftoken.AddArgument(this.curr.Root())
			this.curr = currnode.parent
			scan.Ignore()
			return branchExpressionOperatorPart
		}
		this.curr = currnode
		if scan.IsEOF() {
			this.error("Arguments missing end bracket. End of file reached.")
			return nil
		}
		this.error("Invalid char in Arguments. %c", r)
		return nil
	}
}

//
func branchExpressionOperatorPart(this *parser) stateFn {
	scan := this.scan
	scan.SkipSpaces()

	if scan.IsEOF() {
		return nil
	}
	if this.AcceptOperator() {
		this.parseOperator()
		return branchExpressionValuePart
	}
	if scan.Accept("=") {
		this.parseLRFunc()
		this.curr = this.add(NewGroupToken(""))
		return branchExpressionValuePart
	}
	switch scan.Next() {
	case ')':
		return this.parseCloseBracket()
	}
	scan.Rollback()
	return nil
}
