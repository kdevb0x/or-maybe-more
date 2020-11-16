package net

import (
	"github.com/pion/sdp"
	rtc "github.com/pion/webrtc/v3"
)

var _ = rtc.NewMediaEngine()
var _ = sdp.ConnectionRole
