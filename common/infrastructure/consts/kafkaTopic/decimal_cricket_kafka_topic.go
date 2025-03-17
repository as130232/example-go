package kafkaTopic

const (
	Topic = "decimal.cricket"

	// MatchListTopic 賽事列表 by API
	MatchListTopic = "decimal.cricket.match.list"
	// DeltaTopic 原生賽事資料(賽事、賠率、比分) by websocket
	DeltaTopic = "decimal.cricket.delta"
	// MatchStateTopic 賽事狀態(關賽) by websocket
	MatchStateTopic = "decimal.cricket.match.state"
)
