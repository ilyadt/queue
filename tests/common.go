package queuetest

type ContextKey uint32

const (
	_ ContextKey = iota

	ResultChanContextKey //ctx:value chan *Result
	CancelChanContextKey //ctx:value chan context.CancelFunc
)

type Result struct {
	Num  int    //
	Err  error  //
	Resp string //
}

type ScenarioConfig struct {
	ServerURL string
	QName      string
}

type Scenario struct {
	serverBaseURL string
	queue string
}
