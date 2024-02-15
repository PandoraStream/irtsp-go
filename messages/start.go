package messages

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// StartRequest represents a START request message
type StartRequest struct {
	Message
	Scheme  string
	Timings []uint64
}

func (sreq *StartRequest) parse(data []byte) {
	// TODO - Error check this

	contents := string(data)
	fields := strings.Split(contents, "\r\n")
	headers := sreq.headersToMap(fields[3 : len(fields)-2])

	sreq.IRTSPVersion = fields[0]
	sreq.SequenceNumber, _ = strconv.Atoi(strings.Split(fields[1], "=")[1])

	if scheme, ok := headers["sc"]; ok {
		sreq.Scheme = scheme
	}

	if timingsString, ok := headers["t"]; ok {
		timings := strings.Split(timingsString, ",")

		for _, timingString := range timings {
			timing, err := strconv.ParseUint(timingString, 10, 64)
			if err != nil {
				// TODO - Handle this
				continue
			}

			sreq.Timings = append(sreq.Timings, timing)
		}
	}
}

// Bytes encodes the request into a byte slice
func (sreq *StartRequest) Bytes() []byte {
	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("%s\r\n", sreq.IRTSPVersion))
	buffer.WriteString(fmt.Sprintf("Seq=%d\r\n", sreq.SequenceNumber))
	buffer.WriteString("SET/START\r\n")

	if sreq.Scheme == "" {
		buffer.WriteString("sc\r\n")
	} else {
		buffer.WriteString(fmt.Sprintf("sc=%s\r\n", sreq.Scheme))
	}

	timings := make([]string, len(sreq.Timings))

	for i, timing := range sreq.Timings {
		timings[i] = strconv.Itoa(int(timing))
	}

	buffer.WriteString(fmt.Sprintf("t=%s\r\n", strings.Join(timings, ",")))
	buffer.WriteString("Submit\r\n")

	return buffer.Bytes()
}

// StartResponse represents a START response message
type StartResponse struct {
	Message
	Status  int
	Scheme  string
	Timings []uint64
}

func (srsp *StartResponse) parse(data []byte) {
	// TODO - Error check this

	contents := string(data)
	fields := strings.Split(contents, "\r\n")
	headers := srsp.headersToMap(fields[3 : len(fields)-2])

	srsp.IRTSPVersion = fields[0]
	srsp.Status, _ = strconv.Atoi(strings.Split(fields[2], "/")[2])
	srsp.SequenceNumber, _ = strconv.Atoi(strings.Split(fields[1], "=")[1])

	if timingsString, ok := headers["t"]; ok {
		timings := strings.Split(timingsString, ",")

		for _, timingString := range timings {
			timing, err := strconv.ParseUint(timingString, 10, 64)
			if err != nil {
				// TODO - Handle this
				continue
			}

			srsp.Timings = append(srsp.Timings, timing)
		}
	}

	if scheme, ok := headers["sc"]; ok {
		srsp.Scheme = scheme
	}
}

// Bytes encodes the response into a byte slice
func (srsp *StartResponse) Bytes() []byte {
	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("%s\r\n", srsp.IRTSPVersion))
	buffer.WriteString(fmt.Sprintf("Seq=%d\r\n", srsp.SequenceNumber))
	buffer.WriteString(fmt.Sprintf("RSP/START/%d\r\n", srsp.Status))

	timings := make([]string, len(srsp.Timings))

	for i, timing := range srsp.Timings {
		timings[i] = strconv.Itoa(int(timing))
	}

	buffer.WriteString(fmt.Sprintf("t=%s\r\n", strings.Join(timings, ",")))

	if srsp.Scheme == "" {
		buffer.WriteString("sc\r\n")
	} else {
		buffer.WriteString(fmt.Sprintf("sc=%s\r\n", srsp.Scheme))
	}

	buffer.WriteString("Submit\r\n")

	return buffer.Bytes()
}

// NewStartRequestMessage returns a new START message request
func NewStartRequestMessage(data []byte) *StartRequest {
	request := &StartRequest{
		Timings: make([]uint64, 0),
	}

	request.IRTSPVersion = "iRTSP/1.21"

	if data != nil {
		request.parse(data)
	}

	return request
}

// NewStartResponseMessage returns a new START message response
func NewStartResponseMessage(status int, data []byte) *StartResponse {
	response := &StartResponse{
		Status:  status,
		Timings: make([]uint64, 0),
	}

	response.IRTSPVersion = "iRTSP/1.21"

	if data != nil {
		response.parse(data)
	}

	return response
}
