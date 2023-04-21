package chatgpt

type Wrapper interface {
	CreateSummary(content string) (resp *SummaryResponse, err error)
}
