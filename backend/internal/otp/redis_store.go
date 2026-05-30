package otp

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	keyPrefixRegOTP   = "reg:otp:"
	keyPrefixRegRate  = "reg:otp:rate:"
	keyPrefixLoginOTP = "login:otp:"
	keyPrefixLoginRate = "login:otp:rate:"
)

type PendingRegistration struct {
	OTPHash      string `json:"otp_hash"`
	Attempts     int    `json:"attempts"`
	Fullname     string `json:"fullname"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
}

type PendingLogin struct {
	OTPHash  string `json:"otp_hash"`
	Attempts int    `json:"attempts"`
	Email    string `json:"email"`
	UserID   string `json:"user_id"`
}

type Store struct {
	client *redis.Client
	pepper string
}

func NewStore(client *redis.Client, pepper string) *Store {
	return &Store{client: client, pepper: pepper}
}

func (s *Store) Pepper() string { return s.pepper }

func (s *Store) SavePending(ctx context.Context, email string, pending PendingRegistration, ttl time.Duration) error {
	b, err := json.Marshal(pending)
	if err != nil {
		return err
	}
	return s.client.Set(ctx, keyPrefixRegOTP+email, b, ttl).Err()
}

func (s *Store) GetPending(ctx context.Context, email string) (*PendingRegistration, error) {
	val, err := s.client.Get(ctx, keyPrefixRegOTP+email).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var pending PendingRegistration
	if err := json.Unmarshal(val, &pending); err != nil {
		return nil, err
	}
	return &pending, nil
}

func (s *Store) DeletePending(ctx context.Context, email string) error {
	return s.client.Del(ctx, keyPrefixRegOTP+email).Err()
}

func (s *Store) CheckRegRateLimit(ctx context.Context, email string, max int, window time.Duration) error {
	return s.checkRate(ctx, keyPrefixRegRate+email, max, window)
}

func (s *Store) SaveLogin(ctx context.Context, email string, pending PendingLogin, ttl time.Duration) error {
	b, err := json.Marshal(pending)
	if err != nil {
		return err
	}
	return s.client.Set(ctx, keyPrefixLoginOTP+email, b, ttl).Err()
}

func (s *Store) GetLogin(ctx context.Context, email string) (*PendingLogin, error) {
	val, err := s.client.Get(ctx, keyPrefixLoginOTP+email).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var pending PendingLogin
	if err := json.Unmarshal(val, &pending); err != nil {
		return nil, err
	}
	return &pending, nil
}

func (s *Store) DeleteLogin(ctx context.Context, email string) error {
	return s.client.Del(ctx, keyPrefixLoginOTP+email).Err()
}

func (s *Store) CheckLoginRateLimit(ctx context.Context, email string, max int, window time.Duration) error {
	return s.checkRate(ctx, keyPrefixLoginRate+email, max, window)
}

func (s *Store) checkRate(ctx context.Context, key string, max int, window time.Duration) error {
	count, err := s.client.Incr(ctx, key).Result()
	if err != nil {
		return err
	}
	if count == 1 {
		if err := s.client.Expire(ctx, key, window).Err(); err != nil {
			return err
		}
	}
	if int(count) > max {
		return fmt.Errorf("too many otp requests")
	}
	return nil
}

func NewRedisClient(redisURL, host, port, password string, db int) (*redis.Client, error) {
	if redisURL != "" {
		opts, err := redis.ParseURL(redisURL)
		if err != nil {
			return nil, err
		}
		return redis.NewClient(opts), nil
	}
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       db,
	}), nil
}
