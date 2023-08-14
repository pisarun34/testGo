package database

import "context"

type RedisClientInterface interface {
	GetSeeksterToken(ctx context.Context, ssoid string) (string, error)
	SetSeeksterToken(ctx context.Context, ssoid string, seeksterToken string) error
	CloseRedis()
}
