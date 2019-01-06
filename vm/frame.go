package vm

import "dmeijboom/config/tokens"

type FrameKind int

const (
	RootFrame FrameKind = iota
	FunctionFrame
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

func (frame *Frame) Get(name string) *Value {
	if value, exist := frame.Data[name]; exist {
		return value
	} else if frame.Parent != nil {
		return frame.Parent.Get(name)
	}

	return nil
}
