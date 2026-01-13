package data_splitting

import (
	"fmt"
	"github.com/google/uuid"
	"strings"
)

type TargetMessage struct {
	Message     string `json:"message"` // 分割后的消息内容
	Size        int64  `json:"size"`    // 当前分块的实际字节大小
	Index       int    `json:"index"`   // 分块索引（从1开始）
	TotalChunks int    `json:"total_chunks"`
	Offset      int64  `json:"offset"`
	ChunkId     string `json:"chunk_id"`
}

func split(s string, maxByteSize int64, offset int64) ([]TargetMessage, error) {
	if len(s) == 0 {
		return nil, errEmpty
	}
	if maxByteSize <= 0 {
		return nil, fmt.Errorf("maxByteSize必须大于0")
	}

	// 预分配结果切片，避免多次扩容
	estimatedChunks := (len(s) + int(maxByteSize) - 1) / int(maxByteSize)
	result := make([]TargetMessage, 0, estimatedChunks)

	chunkId := uuid.New().String()

	var currentChunk strings.Builder
	currentSize := 0
	chunkIndex := 1

	for _, r := range s { // 直接遍历字符串，无需转换为rune数组
		charStr := string(r)
		charSize := len(charStr)

		// 如果加入当前字符会超过限制，且当前块不为空
		if int64(currentSize+charSize) > maxByteSize && currentSize > 0 {
			result = append(result, TargetMessage{
				Message:     currentChunk.String(),
				Size:        int64(currentSize),
				Index:       chunkIndex,
				Offset:      offset, // 根据实际需求调整
				ChunkId:     chunkId,
				TotalChunks: 0,
			})

			// 重置当前块
			currentChunk.Reset()
			currentChunk.WriteString(charStr)
			currentSize = charSize
			chunkIndex++
			offset += int64(currentSize) // 更新偏移量
		} else {
			currentChunk.WriteRune(r)
			currentSize += charSize
		}
	}

	// 处理最后一块
	if currentSize > 0 {
		result = append(result, TargetMessage{
			Message:     currentChunk.String(),
			Size:        int64(currentSize),
			Index:       chunkIndex,
			Offset:      offset,
			ChunkId:     chunkId,
			TotalChunks: 0,
		})
	}

	// 批量更新总分块数
	totalChunks := len(result)
	for i := range result {
		result[i].TotalChunks = totalChunks
	}

	return result, nil
}
