package service

import "context"

type Creator[Req any, Resp any] interface {
	Create(ctx context.Context, req *Req) (*Resp, error)
	Validate(ctx context.Context, req *Req) error
}

type Reader[Req any, Resp any] interface {
	Read(ctx context.Context, req *Req) (*Resp, error)
	Validate(ctx context.Context, req *Req) error
}
