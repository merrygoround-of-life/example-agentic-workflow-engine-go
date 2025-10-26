package domain

// Request 요청 구조체
type Request struct {
	Message string `json:"message"`
}

// ResponseChunk 에이전트 응답 스트리밍 구조체
type ResponseChunk struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

// FanInChunk 채널 병합용 구조체
type FanInChunk struct {
	Order int           // 에이전트 인덱스
	Chunk ResponseChunk // 실제 데이터
	Done  bool          // 완료 신호
}

// Agent 에이전트 인터페이스
type Agent interface {
	Run(input Request) <-chan ResponseChunk
}
