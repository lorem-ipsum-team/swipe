package domain

type Swipe struct {
	Init       UserID
	Target     UserID
	InitResp   *bool
	TargetResp *bool
}
