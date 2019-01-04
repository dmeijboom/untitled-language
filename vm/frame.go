package vm

import "dmeijboom/config/tokens"

type FrameKind int

const (
	FunctionFrame FrameKind = iota
	BlockFrame
)

type Frame struct {
	Kind FrameKind
	Parent *Frame
	FunctionName string
	Data map[string]*Value
	Types map[string]*Type
	Location *tokens.Location
}

func NewFrame(kind FrameKind, loc *tokens.Location) *Frame {
	return &Frame{
		Kind: kind,
		Location: loc,
		Data: map[string]*Value{},
		Types: map[string]*Type{},
	}
}
