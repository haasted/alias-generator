package impl

const (
	ThisIsAConstant = "Const value"
	AnotherConstant = "Yet more"
	privateConstant = "Wow, they all align"
	_               = "This is a pointless const, that should not be aliased"
)

var (
	ThisIsAVar = "Variable value"
	AnotherVar = "again..."
	andAnother = "..."
	// noalias
	MoreVars = "...."
)

type CodeStruct struct{}

func (CodeStruct) Method1() {}

func NewCode() CodeStruct {
	return CodeStruct{}
}

type _ struct {
	WhatsThePoint string
}
