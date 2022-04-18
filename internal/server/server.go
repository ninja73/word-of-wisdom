package server

import (
	"bufio"
	"errors"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/zeebo/xxh3"
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

const (
	defaultBitStrength = 20
	defaultTimeout     = 5 * time.Second
)

type server struct {
	bitStrength int32
	timeout     time.Duration
	quoteStore  store.Store
	mCache      cache.Cache
	secretKey   string
	bufPool     *bufPool
	ttl         int64
}

func NewServer(quoteStore store.Store, mCache cache.Cache, opts ...Option) *server {
	srv := &server{
		bitStrength: defaultBitStrength,
		timeout:     defaultTimeout,
		quoteStore:  quoteStore,
		mCache:      mCache,
		bufPool:     &bufPool{},
	}

	for _, opt := range opts {
		opt(srv)
	}

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

func (s *server) createSignature(uid string, timestamp int64) uint64 {
	return xxh3.HashString(uid + strconv.FormatInt(int64(s.bitStrength), 10) + strconv.FormatInt(timestamp, 10) + s.secretKey)
}

func (s *server) validateChallenge(challenge *dto.Challenge) error {
	signature := s.createSignature(challenge.Uid, challenge.Timestamp)
	if challenge.Signature != signature {
		return signatureError
	}

	now := time.Now().Unix()
	if challenge.Timestamp+s.ttl < now {
		return expiredError
	}

	hc := pow.NewHashCash(s.bitStrength, challenge.Uid, challenge.Timestamp, challenge.Signature)
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

	if ok := s.mCache.ContainsOrAdd(challenge.Signature); ok {
		s.responseError(conn, challengeExistError)
		return
	}

	s.responseQuote(conn)
}

func (s *server) responseChallenge(conn net.Conn) {
	uid := uuid.New().String()
	timeout := time.Now().Unix()
	signature := s.createSignature(uid, timeout)

	challenge := &dto.Challenge{
		BitStrength: s.bitStrength,
		Uid:         uid,
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
	quote := s.quoteStore.RandomQuote()
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

	if err := conn.SetDeadline(time.Now().Add(s.timeout)); err != nil {
		log.Error(err)
		s.responseError(conn, internalServerError)
		return
	}

	reader := bufio.NewReader(conn)

	data := s.bufPool.get()
	defer s.bufPool.put(data)

	size, err := reader.Read(data)
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

	case dto.Type_REQUEST_CHALLENGE:

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

	for {
		conn, err := listener.Accept()
		if err != nil {

			continue
		}
		go s.handlerConn(conn)
	}
}
