package protocol

// MessageBuffer implements a fixed capacity ringbuffer for items of type
// MessageToReportGeneration
type MessageBuffer struct {
	start  int
	length int
	buffer []*MessageToReportGeneration
}

func NewMessageBuffer(cap int) *MessageBuffer {
	return &MessageBuffer{
		0,
		0,
		make([]*MessageToReportGeneration, cap),
	}
}

// Peek at the front item
func (rb *MessageBuffer) Peek() *MessageToReportGeneration {
	if rb.length == 0 {
		return nil
	} else {
		return rb.buffer[rb.start]
	}
}

// Pop front item
func (rb *MessageBuffer) Pop() *MessageToReportGeneration {
	result := rb.Peek()
	if result != nil {
		rb.buffer[rb.start] = nil
		rb.start = (rb.start + 1) % len(rb.buffer)
		rb.length--
	}
	return result
}

// Push new item to back. If the additional item would lead
// to the capacity being exceeded, remove the front item first
func (rb *MessageBuffer) Push(msg MessageToReportGeneration) {
	if rb.length == len(rb.buffer) {
		rb.Pop()
	}
	rb.buffer[(rb.start+rb.length)%len(rb.buffer)] = &msg
	rb.length++
}
