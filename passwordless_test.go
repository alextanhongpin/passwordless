package passwordless_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/alextanhongpin/passwordless"
)

func TestPasswordlessStart(t *testing.T) {
	var (
		userID = "xyz"
	)
	t.Run("when create success", func(t *testing.T) {
		repo := &mockCodeRepository{
			createCode: passwordless.NewCode(),
			createErr:  nil,
		}
		svc := passwordless.New(repo)
		res, err := svc.Start(context.Background(), userID)
		if err != nil {
			t.Fatal(err)
		}
		if res.Code == "" {
			t.Fatalf("expected .Code to have value, actual is empty")
		}
	})

	t.Run("when create fail", func(t *testing.T) {
		repo := &mockCodeRepository{
			createErr: sql.ErrNoRows,
		}
		svc := passwordless.New(repo)
		res, err := svc.Start(context.Background(), userID)
		if !errors.Is(err, sql.ErrNoRows) {
			t.Fatalf("expected %s, got %s", sql.ErrNoRows, err)
		}
		if res != nil {
			t.Fatalf("expected response to be nil, got %v", res)
		}
	})
}

func TestPasswordlessAuthorize(t *testing.T) {

	var (
		code = "xyz"
	)
	t.Run("when code does not exists", func(t *testing.T) {
		repo := &mockCodeRepository{
			findErr: sql.ErrNoRows,
		}
		svc := passwordless.New(repo)
		res, err := svc.Authorize(context.Background(), code)
		if !errors.Is(sql.ErrNoRows, err) {
			t.Fatalf("expected %v, got %v", sql.ErrNoRows, err)
		}
		if res != nil {
			t.Fatalf("expected response to be nil, got %v", res)
		}
	})

	t.Run("when code expired", func(t *testing.T) {
		codeInDB := passwordless.NewCode()
		codeInDB.CreatedAt = time.Now().Add(-15 * time.Minute)
		repo := &mockCodeRepository{
			findCode: codeInDB,
		}
		svc := passwordless.New(repo)
		res, err := svc.Authorize(context.Background(), codeInDB.Code)
		if !errors.Is(passwordless.ErrCodeExpired, err) {
			t.Fatalf("expected %v, got %v", passwordless.ErrCodeExpired, err)
		}
		if res != nil {
			t.Fatalf("expected response to be nil, got %v", res)
		}
	})

	t.Run("when delete fail", func(t *testing.T) {
		codeInDB := passwordless.NewCode()
		repo := &mockCodeRepository{
			findCode:  codeInDB,
			deleteErr: sql.ErrNoRows,
		}
		svc := passwordless.New(repo)
		res, err := svc.Authorize(context.Background(), codeInDB.Code)
		if !errors.Is(sql.ErrNoRows, err) {
			t.Fatalf("expected %v, got %v", passwordless.ErrCodeExpired, err)
		}
		if res != nil {
			t.Fatalf("expected response to be nil, got %v", res)
		}
	})
}

type mockCodeRepository struct {
	createCode *passwordless.Code
	createErr  error
	deleteErr  error
	findCode   *passwordless.Code
	findErr    error
}

func (m *mockCodeRepository) Create(ctx context.Context, userID string) (*passwordless.Code, error) {
	return m.createCode, m.createErr
}

func (m *mockCodeRepository) Delete(ctx context.Context, code string) error {
	return m.deleteErr
}

func (m *mockCodeRepository) Find(ctx context.Context, code string) (*passwordless.Code, error) {
	return m.findCode, m.findErr
}
