package server

import (
	"bufio"
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/zeebo/xxh3"
	"golang.org/x/time/rate"
	"google.golang.org/protobuf/proto"
	"net"
	"strconv"
	"time"
	"wow/internal/cache"
	"wow/internal/store"
	"wow/package/pow"
	"wow/package/proto/dto"
)

var (
	internalServerError = errors.New("internal server error")
	badRequestError     = errors.New("bad request")
	signatureError      = errors.New("signature not valid")
	expiredError        = errors.New("expired error")
	proofError          = errors.New("challenge proof error")
	challengeExistError = errors.New("challenge already exist")
)

type server struct {
	quoteStore store.Store
	powCache   cache.Cache
	reqPool    *reqPool
	Limiter    *rate.Limiter
	*Options
}

func NewServer(quoteStore store.Store, powCache cache.Cache, opts ...*Options) *server {
	srv := &server{
		quoteStore: quoteStore,
		powCache:   powCache,
		reqPool:    new(reqPool),
	}

	if len(opts) != 0 {
		srv.Options = opts[0]
	} else {
		srv.Options = &Options{}
	}

	setDefaultWorkOptions(srv.Options)

	srv.Limiter = rate.NewLimiter(rate.Limit(srv.Limit), int(srv.Limit+(srv.Limit*10)/100))

	return srv
}

func (s *server) writeResponse(conn net.Conn, msg *dto.Msg) {
	resp, err := proto.Marshal(msg)
	if err != nil {
		log.Error(err)
		s.responseError(conn, internalServerError)
		return
	}

	if _, err := conn.Write(resp); err != nil {
		log.Error(err)
	}
}

func (s *server) createSignature(bitStrength int32, data string, timestamp int64) uint64 {
	return xxh3.HashString(data + strconv.FormatInt(int64(bitStrength), 10) + strconv.FormatInt(timestamp, 10) + s.SecretKey)
}

func (s *server) validateChallenge(challenge *dto.Challenge) error {
	signature := s.createSignature(s.BitStrength, challenge.Data, challenge.Timestamp)
	if challenge.Signature != signature {
		return signatureError
	}

	now := time.Now()
	ttl := time.Unix(challenge.Timestamp, 0).Add(s.Expiration)

	if now.After(ttl) {
		return expiredError
	}

	hc := pow.HashCash{
		BitStrength: s.BitStrength,
		Data:        challenge.Data,
		Timestamp:   challenge.Timestamp,
		Counter:     challenge.Counter,
		Signature:   challenge.Signature,
	}

	if !hc.Check() {
		return proofError
	}

	return nil
}

func (s *server) responseProofChallenge(conn net.Conn, data []byte) {
	challenge := new(dto.Challenge)
	if err := proto.Unmarshal(data, challenge); err != nil {
		log.Error(err)
		s.responseError(conn, internalServerError)
		return
	}

	if err := s.validateChallenge(challenge); err != nil {
		log.Error(err)
		s.responseError(conn, err)
		return
	}

	ctx := context.Background()

	ok, err := s.powCache.ContainsOrAdd(ctx, challenge.Signature)
	if err != nil {
		log.Error(err)
		s.responseError(conn, internalServerError)
		return
	}

	if ok {
		s.responseError(conn, challengeExistError)
		return
	}

	s.responseQuote(conn)
}

func (s *server) responseChallenge(conn net.Conn) {
	rnd := randomString(16)
	timeout := time.Now().Unix()
	signature := s.createSignature(s.BitStrength, rnd, timeout)

	challenge := &dto.Challenge{
		BitStrength: s.BitStrength,
		Data:        rnd,
		Timestamp:   timeout,
		Signature:   signature,
	}

	data, err := proto.Marshal(challenge)
	if err != nil {
		log.Error(err)
		s.responseError(conn, internalServerError)
		return
	}

	msg := &dto.Msg{
		Type: dto.Type_RESPONSE_CHALLENGE,
		Data: data,
	}

	s.writeResponse(conn, msg)
}

func (s *server) responseQuote(conn net.Conn) {
	ctx := context.Background()
	quote, err := s.quoteStore.RandomQuote(ctx)
	if err != nil {
		log.Error(err)
		s.responseError(conn, internalServerError)
		return
	}

	msg := &dto.Msg{
		Type: dto.Type_RESPONSE_QUOTE,
		Data: []byte(quote),
	}

	s.writeResponse(conn, msg)
}

func (s *server) responseError(conn net.Conn, err error) {
	msg := &dto.Msg{
		Type: dto.Type_RESPONSE_ERROR,
		Data: []byte(err.Error()),
	}

	s.writeResponse(conn, msg)
}

func (s *server) handlerConn(conn net.Conn) {
	defer conn.Close()

	if err := conn.SetDeadline(time.Now().Add(s.Timeout)); err != nil {
		log.Error(err)
		s.responseError(conn, internalServerError)
		return
	}

	data := s.reqPool.get()
	defer s.reqPool.put(data)

	size, err := bufio.NewReader(conn).Read(data)
	if err != nil {
		log.Error(err)
		s.responseError(conn, internalServerError)
		return
	}

	msg := &dto.Msg{}
	if err := proto.Unmarshal(data[:size], msg); err != nil {
		log.Error(err)
		s.responseError(conn, badRequestError)
		return
	}

	switch msg.Type {
	case dto.Type_REQUEST_QUOTE:
		if s.Limiter.Allow() {
			s.responseQuote(conn)
			return
		}

		s.responseChallenge(conn)
	case dto.Type_REQUEST_CHALLENGE:
		s.responseProofChallenge(conn, msg.Data)
	default:
		s.responseError(conn, badRequestError)
	}
}

func (s *server) TCPListen(address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	defer listener.Close()

	log.Infof("Server started: %s", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Error(err)
			continue
		}
		go s.handlerConn(conn)
	}
}
