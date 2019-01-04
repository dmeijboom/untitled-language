package ast

type Visitor interface {
	VisitIdent(ident *Ident)
	VisitField(field *Field)
	VisitType(type_ *Type)
	VisitPreInitialize(init *Initialize)
	VisitInitialize(init *Initialize)
	VisitInitializeField(initField *InitializeField)
	VisitLiteral(literal *Literal)
	VisitBlock(block *Block)
	VisitSection(section *Section)
	VisitTypedef(typedef *Typedef)
	VisitAssign(assign *Assign)
	VisitSource(source *Source)
}

type Visitable interface {
	Accept(visitor Visitor)
}
