package eventStatus

const (
	MatchScheduled    = 0  // 比賽尚未開始
	MatchStable       = 1  // 比賽進行中
	BallInProgress    = 2  // 比賽進行中
	UnknownEvent      = 3  // 當無法立即得知結果時會使用
	RunsLikely        = 4  // 可能邊界球
	PossibleWicket    = 5  // 可能 wicket
	UmpireReview      = 6  // 複查
	BreakInPlay       = 7  // 中斷
	DataAssemblyIssue = 8  // 需要更正分數或來源暫時無法使用時發送
	SuperOver         = 9  // 打平後加比
	MatchComplete     = 10 // 比賽結束
	EventClosed       = 11 // 事件結束
	// 12 ?
	Ejected = 13 // 當某事件強制從清單刪除時會發送並附加 EventClosed 狀態
)
