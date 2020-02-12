package db

// This is a thoroughly documented
// function.
// noalias
// Should it be exempted from aliasing?
func CreateDatabase() {}

// noalias
type Database struct{}

// noalias
type (
	Connection struct{}
	Index      struct{}
)

// noalias
var (
	Variable1 = false
	Variable2 = 6
)
