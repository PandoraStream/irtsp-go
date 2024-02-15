// Package messages implements iRTSP message structs
package messages

import "strings"

// Message represents a generic message. This struct is embedded in all other message structs
type Message struct {
	IRTSPVersion   string
	SequenceNumber int
}

func (m *Message) headersToMap(headers []string) map[string]string {
	headersMap := make(map[string]string)

	for _, header := range headers {
		parts := strings.Split(header, "=")
		key := parts[0]
		value := ""

		if len(parts) > 1 {
			value = parts[1]
		}

		headersMap[key] = value
	}

	return headersMap
}
