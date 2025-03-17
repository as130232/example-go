package apiUrl

const (
	// 呼叫方式 header X-API-KEY: <client API key>
	WssMatchFeedStg  = "wss://staging-delta-feed.decimalcricket.net/api/match/feed"
	WssMatchFeedProd = "wss://delta-feed.decimalcricket.net/api/match/feed"

	// 呼叫方式 curl --header "X-API-KEY: <client API key>" https://staging-fixtures.decimalcricket.net/api/match/list
	HttpMatchListStg  = "https://staging-fixtures.decimalcricket.net/api/match/list"
	HttpMatchListProd = "https://fixtures.decimalcricket.net/api/match/list"
)
