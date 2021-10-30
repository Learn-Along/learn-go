package codegenerator

// AstPrinter is a Visitor implementation to print Abstract Syntax Tree (AST) in lisp style
type AstPrinter struct {

}

func (a *AstPrinter) visitUnionExpr(u UnionExpression) Output {
	return nil
}

func (a *AstPrinter) visitSelectExpr(s SelectExpression) Output {
	return nil
}

