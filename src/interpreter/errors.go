package interpreter

type BreakError struct{}
type ContinueError struct{}

func (e BreakError) Error() string    { return "break statement" }
func (e ContinueError) Error() string { return "continue statement" }
