package Strings

import "strings"

func SplitText(text string, chunkSize int) []string {
	sentences := strings.Split(text, "。")
	var (
		chunks       []string
		currentChunk strings.Builder
	)
	for _, sentence := range sentences {
		if sentence == "" {
			continue
		}
		if currentChunk.Len()+len(sentence) > chunkSize && currentChunk.Len() > 0 {
			chunks = append(chunks, currentChunk.String()+"。")
			currentChunk.Reset()
		}
		currentChunk.WriteString(sentence)
		currentChunk.WriteString("。")
	}
	if currentChunk.Len() > 0 {
		chunks = append(chunks, currentChunk.String())
	}
	return chunks
}
