package interpreter

type BreakError struct{}
type ContinueError struct{}

func (BreakError) Error() string { return "no enclosing block out of which to break" }
func (ContinueError) Error() string {
	return "no enclosing loop out of which to continue"
}
