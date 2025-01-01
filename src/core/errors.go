package core

type BreakError struct{}

func (BreakError) Error() string { return "no enclosing block out of which to break" }
func NewBreakError() BreakError  { return BreakError{} }

type ContinueError struct{}

func (ContinueError) Error() string   { return "no enclosing loop out of which to continue" }
func NewContinueError() ContinueError { return ContinueError{} }

type ReturnError struct {
	value Value
}

func (ReturnError) Error() string            { return "return statement outside of function body" }
func NewReturnError(value Value) ReturnError { return ReturnError{value} }
func (r ReturnError) Value() Value           { return r.value }