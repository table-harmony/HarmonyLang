package interpreter

type BreakError struct{}
type ContinueError struct{}

func (BreakError) Error() string { return "no enclosing block out of which to break" }
func (ContinueError) Error() string {
	return "no enclosing loop out of which to continue"
}

type ReturnError struct {
	Value Value
}

func (ReturnError) Error() string {
	return "return statement outside of function body"
}
