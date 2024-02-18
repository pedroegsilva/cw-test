package parser

import "fmt"

type SyntaxError struct {
	header  LogHeader
	message string
}

func (e *SyntaxError) Error() string {
	return fmt.Sprintf("[syntax error] %s: %s ", e.header.String(), e.message)
}

type ContextError struct {
	header  LogHeader
	message string
}

func (e *ContextError) Error() string {
	return fmt.Sprintf("[context error] %s: %s ", e.header.String(), e.message)
}
