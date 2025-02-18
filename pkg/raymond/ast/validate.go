package ast

// validateVisitor implements the Visitor interface to print a AST.
type validateVisitor struct {
	simple bool
}

func newValidateVisitor() *validateVisitor {
	return &validateVisitor{
		simple: true,
	}
}

// Validate returns true if given AST is simple, otherwise false.
func Validate(node Node) bool {
	visitor := newValidateVisitor()
	node.Accept(visitor)
	return visitor.simple
}

//
// Visitor interface
//

// Statements

// VisitProgram implements corresponding Visitor interface method
func (v *validateVisitor) VisitProgram(node *Program) interface{} {
	for _, n := range node.Body {
		n.Accept(v)
	}

	return nil
}

// VisitMustache implements corresponding Visitor interface method
func (v *validateVisitor) VisitMustache(node *MustacheStatement) interface{} {
	node.Expression.Accept(v)

	return nil
}

// VisitBlock implements corresponding Visitor interface method
func (v *validateVisitor) VisitBlock(node *BlockStatement) interface{} {
	v.simple = false

	node.Expression.Accept(v)

	if node.Program != nil {
		node.Program.Accept(v)
	}

	if node.Inverse != nil {
		node.Inverse.Accept(v)
	}

	return nil
}

// VisitPartial implements corresponding Visitor interface method
func (v *validateVisitor) VisitPartial(node *PartialStatement) interface{} {
	node.Name.Accept(v)

	if len(node.Params) > 0 {
		node.Params[0].Accept(v)
	}

	// hash
	if node.Hash != nil {
		node.Hash.Accept(v)
	}

	return nil
}

// VisitContent implements corresponding Visitor interface method
func (v *validateVisitor) VisitContent(node *ContentStatement) interface{} {
	return nil
}

// VisitComment implements corresponding Visitor interface method
func (v *validateVisitor) VisitComment(node *CommentStatement) interface{} {
	return nil
}

// Expressions

// VisitExpression implements corresponding Visitor interface method
func (v *validateVisitor) VisitExpression(node *Expression) interface{} {
	// path
	node.Path.Accept(v)

	// params
	for _, n := range node.Params {
		n.Accept(v)
	}

	// hash
	if node.Hash != nil {
		node.Hash.Accept(v)
	}

	return nil
}

// VisitSubExpression implements corresponding Visitor interface method
func (v *validateVisitor) VisitSubExpression(node *SubExpression) interface{} {
	node.Expression.Accept(v)

	return nil
}

// VisitPath implements corresponding Visitor interface method
func (v *validateVisitor) VisitPath(node *PathExpression) interface{} {
	return nil
}

// Literals

// VisitString implements corresponding Visitor interface method
func (v *validateVisitor) VisitString(node *StringLiteral) interface{} {
	return nil
}

// VisitBoolean implements corresponding Visitor interface method
func (v *validateVisitor) VisitBoolean(node *BooleanLiteral) interface{} {
	return nil
}

// VisitNumber implements corresponding Visitor interface method
func (v *validateVisitor) VisitNumber(node *NumberLiteral) interface{} {
	return nil
}

// Miscellaneous

// VisitHash implements corresponding Visitor interface method
func (v *validateVisitor) VisitHash(node *Hash) interface{} {
	for _, p := range node.Pairs {
		p.Accept(v)
	}
	return nil
}

// VisitHashPair implements corresponding Visitor interface method
func (v *validateVisitor) VisitHashPair(node *HashPair) interface{} {
	node.Val.Accept(v)

	return nil
}
