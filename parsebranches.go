package goexpression

func (this *parser) pumpExpression() {
	this.state = this.branchExpressionValuePart
	for this.state != nil {
		this.state()
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
func (this *parser) branchExpressionValuePart() bool {
	scan := this.scan
	scan.SkipSpaces()
	if scan.IsEOF() {
		this.state = nil
		return true
	}
	if scan.ScanNumber() {
		this.state = this.branchExpressionOperatorPart
		this.add(NewNumberToken(scan.Commit()))
		return true
	}
	if scan.ScanWord() {
		this.state = this.branchExpressionOperatorPart
		this.add(NewIdentityToken(scan.Commit()))
		return true
	}
	switch scan.Next() {
	case '(':
		this.state = this.branchExpressionValuePart
		this.parseOpenBracket()
		return true
	}
	this.error("Unexpected token. ")
	this.state = nil
	return false
}

//
func (this *parser) branchExpressionOperatorPart() bool {
	scan := this.scan
	scan.SkipSpaces()

	if scan.IsEOF() {
		this.state = nil
		return true
	}
	if scan.Accept("+-*/") {
		this.state = this.branchExpressionValuePart
		this.parseOperator()
		return true
	}
	if scan.Accept("=") {
		this.state = this.branchExpressionValuePart
		this.parseLRFunc()
		this.curr = this.add(NewGroupToken(""))
		return true
	}
	switch scan.Next() {
	case ')':
		this.state = this.branchExpressionOperatorPart
		this.parseCloseBracket()
		return true
	}
	scan.Rollback()
	this.state = nil
	return false
}
