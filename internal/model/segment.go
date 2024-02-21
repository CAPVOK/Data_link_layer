package model

type Segment struct {
	Sender           string `json:"sender"`
	Timestamp        string `json:"timestamp"`
	AmountOfSegments int    `json:"amountOfSegments"`
	SegmentNum       int    `json:"segmentNum"`
	Message          string `json:"message"`
}
