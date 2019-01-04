package vm

type CallStack struct {
	frames []*Frame
}

func NewCallStack() *CallStack {
	return &CallStack{frames: []*Frame{}}
}

func (callStack *CallStack) Push(frame *Frame) {
	frame.Parent = callStack.Frame()
	callStack.frames = append(callStack.frames, frame)
}

func (callStack *CallStack) Pop() *Frame {
	frame := callStack.Frame()
	callStack.frames = callStack.frames[:len(callStack.frames)-1]
	return frame
}

func (callStack *CallStack) Frame() *Frame {
	if len(callStack.frames) == 0 {
		return nil
	}

	return callStack.frames[len(callStack.frames)-1]
}
