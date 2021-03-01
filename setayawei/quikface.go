package page

import (
	"bytes"
	"context"
	"crypto"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/crypto/nacl/box"

	rtc "github.com/pion/webrtc/v3"
	uuid "github.com/satori/go.uuid"
)

// QfSession is the QuikFace Aplication Programming Interface
type QfSession struct {
	URL     *url.URL
	Created time.Time
	// time in unixnano
	expiration int64

	Duration time.Duration
	Key      crypto.PublicKey

	LastErr error

	// the number of active participants in the current sesssion.
	UserCount int

	// upper bound for participants; -1 == unlimited.
	UserLimit int

	// unexported local field context for the session
	ctx context.Context

	// scrypt hash for join key comparison
	JoinHash []byte

	TLSConfig *tls.Config
}

// Session implements QfAPI.
func (qfs *QfSession) Session() *QfSession {
	return qfs
}

// ServeHTTP implements http.Handler interface.
// This handler contains the session specific operations, such as those relating
// to webrtc. At the point this handler is called, the session is verified to
// exist, and the client related to r has been identidied and authenticated.
func (qfs *QfSession) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// client should encrypt the joinkey
	// sent in reply from the GET response
	// with thier private key, and send that
	// value in the POST, along with their
	// public key, to be used for webrtc
	// session signaling.

	// if all is well, pass on to the
	// handler:
	// sesh.ServeHTTP(w, r)

	// get stream reader for webrtc signaling
	var signalbuff = new(bytes.Buffer)

	// because bytes.Buffer.ReadFrom() can
	// panic if the buffer grows too large,
	// we need to catch any panics without
	// disrupting execution, so we do it in
	// a function litteral.
	catchOverflow := func(rd io.Reader, e error) {
		defer recover()
		_, err := signalbuff.ReadFrom(rd)
		if err != nil {
			e = fmt.Errorf("%w\n", err)
			return
		}
	}

	if sr, err := r.MultipartReader(); err == nil {
		var n *multipart.Part
		var err error
		for {
			n, err = sr.NextPart()
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
			}

			catchOverflow(n, err)
		}

	}

	// TODO: create/join webrtc session
	var _ rtc.API

}

var webrtcAPI rtc.API
var webrtcAPIConfig rtc.Configuration

func webrtcSession() {
	webrtcAPIConfig = rtc.Configuration{
		ICEServers: []rtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	m := new(rtc.MediaEngine)

	if err := m.RegisterCodec(rtc.RTPCodecParameters{
		RTPCodecCapability: rtc.RTPCodecCapability{MimeType: "video/VP8", ClockRate: 9000, Channels: 0, SDPFmtpLine: "", RTCPFeedback: nil},
		PayloadType:        96,
	}, rtc.RTPCodecTypeVideo); err != nil {
		panic(err)
	}

	if err := m.RegisterCodec(rtc.RTPCodecParameters{
		RTPCodecCapability: rtc.RTPCodecCapability{MimeType: "audio/opus", ClockRate: 48000, Channels: 0, SDPFmtpLine: "", RTCPFeedback: nil},
		PayloadType:        111,
	}, rtc.RTPCodecTypeAudio); err != nil {
		panic(err)
	}

	webrtcAPI = *rtc.NewAPI(rtc.WithMediaEngine(m))
}

func (qfs *QfSession) Join(joinkey []byte, clientPubkey []byte) (http.Handler, error) {
	var _ = box.Overhead
}

func (qfs *QfSession) Valid(current time.Time) bool {
	if current.UnixNano() < qfs.expiration {
		return true
	}
	return false
}

type QfAPI interface {
	Session() *QfSession
}

func (s *Server) NewQfSession(ctx context.Context, userlimit int) *QfSession {
	ses := new(QfSession)
	ses.Created = time.Now()
	ses.UserLimit = userlimit
	if ctx != nil {
		ses.ctx = context.WithValue(ctx, nil, nil)
	}
	id := uuid.NewV4()
	// TODO: build url and make handler

	s.qfsessions = append(s.qfsessions, ses)
	s.QfSessions[id] = uint64(len(s.qfsessions) - 1)

	// BUG:
	// copies the servers tls config to the session, should it create its
	// own instead of cloning?
	ses.TLSConfig = s.TLSConfig.Clone()
	return ses
}
