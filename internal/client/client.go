package client

import (
	"bufio"
	"errors"
	"fmt"
	"google.golang.org/protobuf/proto"
	"net"
	"time"
	"wow/package/pow"
	"wow/package/proto/dto"
)

type client struct {
	serverAddress string
	timeout       time.Duration
}

func NewClient(serverAddress string, timeout time.Duration) *client {
	return &client{
		serverAddress: serverAddress,
		timeout:       timeout,
	}
}

func (c *client) GetQuote() (string, error) {
	return c.quote(&dto.Msg{Type: dto.Type_REQUEST_QUOTE, Data: make([]byte, 10)})
}

func (c *client) quote(quoteMsg *dto.Msg) (string, error) {
	resp, err := c.request(quoteMsg)
	if err != nil {
		return "", err
	}

	switch resp.Type {
	case dto.Type_RESPONSE_QUOTE:
		return string(resp.Data), err
	case dto.Type_RESPONSE_CHALLENGE:
		msg, err := c.challenge(resp.Data)
		if err != nil {
			return "", err
		}
		return c.quote(msg)
	case dto.Type_RESPONSE_ERROR:
		return "", errors.New(string(resp.Data))
	default:
		return "", fmt.Errorf("unknown message %s", resp.Type)
	}
}

func (c *client) request(msg *dto.Msg) (*dto.Msg, error) {
	conn, err := net.DialTimeout("tcp", c.serverAddress, c.timeout)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	req, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}

	if err := conn.SetDeadline(time.Now().Add(c.timeout)); err != nil {
		return nil, err
	}

	if _, err := conn.Write(req); err != nil {
		return nil, err
	}

	data := make([]byte, 4096)
	size, err := bufio.NewReader(conn).Read(data)
	if err != nil {
		return nil, err
	}

	var resp dto.Msg
	if err := proto.Unmarshal(data[:size], &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *client) challenge(data []byte) (*dto.Msg, error) {
	var challenge dto.Challenge
	if err := proto.Unmarshal(data, &challenge); err != nil {
		return nil, err
	}

	hc := pow.HashCash{
		BitStrength: challenge.BitStrength,
		Data:        challenge.Data,
		Timestamp:   challenge.Timestamp,
		Signature:   challenge.Signature,
	}
	hc.FindProof()

	challenge.Counter = hc.Counter
	challengeData, err := proto.Marshal(&challenge)
	if err != nil {
		return nil, err
	}

	msg := &dto.Msg{
		Type: dto.Type_REQUEST_CHALLENGE,
		Data: challengeData,
	}

	return msg, nil
}
