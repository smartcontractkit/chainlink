package msgbuf

// MessageBuffer implements a fixed capacity ringbuffer for items of type
// []byte.
type MessageBuffer struct {
	start  int
	length int
	buffer [][]byte
}

func NewMessageBuffer(cap int) *MessageBuffer {
	return &MessageBuffer{
		0,
		0,
		make([][]byte, cap),
	}
}

// Peek at the front item
func (rb *MessageBuffer) Peek() []byte {
	if rb.length == 0 {
		return nil
	} else {
		return rb.buffer[rb.start]
	}
}

// Pop front item
func (rb *MessageBuffer) Pop() []byte {
	result := rb.Peek()
	if result != nil {
		rb.buffer[rb.start] = nil
		rb.start = (rb.start + 1) % len(rb.buffer)
		rb.length--
	}
	return result
}

// Push new item to back. If the additional item would lead to the capacity
// being exceeded, remove the front item first.
//
// Returns the removed front item, or nil.
func (rb *MessageBuffer) Push(msg []byte) []byte {
	var result []byte

	if msg == nil {
		panic("cannot push nil")
	}
	if rb.length == len(rb.buffer) {
		result = rb.Pop()
	}
	rb.buffer[(rb.start+rb.length)%len(rb.buffer)] = msg
	rb.length++
	return result
}
