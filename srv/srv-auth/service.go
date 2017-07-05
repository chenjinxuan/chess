package main

import (
	"errors"
	"regexp"
	"golang.org/x/net/context"
	. "chess/srv/srv-auth/proto"
	"chess/common/log"
	"time"
)

var (
	ERROR_METHOD_NOT_SUPPORTED = errors.New("method not supoorted")
)

var (
	uuid_regexp = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
)

type server struct {
}

func (s *server) init() {
}

func (s *server) Auth(ctx context.Context, args *AuthArgs) (*AuthRes, error) {
	// TODO check token
	log.Debugf("AuthArgs(%+v)", *args)
	return &AuthRes{Ret: 1, Msg: "ok"}, nil
}

func (s *server) TokenInfo(ctx context.Context, args *TokenInfoArgs) (*TokenInfoRes, error) {
	log.Debugf("TokenInfoArgs(%+v)", *args)

	return &TokenInfoRes{Expire: time.Now().Unix()}, nil
}

func (s *server) RefreshToken(ctx context.Context, args *RefreshTokenArgs) (*RefreshTokenRes, error) {
	log.Debugf("RefreshTokenArgs(%+v)", *args)

	return &RefreshTokenRes{Ret: 1, Msg: "ok"}, nil
}

func (s *server) DestroyToken(ctx context.Context, args *DestroyTokenArgs) (*DestroyTokenRes, error) {
	log.Debugf("RefreshTokenArgs(%+v)", *args)

	return &DestroyTokenRes{Ret: 1}, nil
}