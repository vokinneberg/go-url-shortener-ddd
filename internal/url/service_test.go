package url

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
)

func computeShort(original string) string {
	h := sha1.New()
	h.Write([]byte(original))
	return hex.EncodeToString(h.Sum(nil))[:8]
}

func TestURLService_Shorten_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockURLReaderWriter(ctrl)
	svc := NewURLService(mockRepo)

	original := "https://example.com/a"
	expectedID := computeShort(original)

	mockRepo.
		EXPECT().
		Save(gomock.Any()).
		DoAndReturn(func(u *URL) error {
			if u == nil {
				t.Fatalf("expected non-nil URL")
			}
			if u.ID != expectedID {
				t.Fatalf("expected ID %q, got %q", expectedID, u.ID)
			}
			if u.Original != original {
				t.Fatalf("expected Original %q, got %q", original, u.Original)
			}
			return nil
		}).
		Times(1)

	got, err := svc.Shorten(original)
	if err != nil {
		t.Fatalf("Shorten() unexpected error: %v", err)
	}
	if got == nil {
		t.Fatalf("Shorten() returned nil url")
	}
	if got.ID != expectedID {
		t.Fatalf("Shorten() expected ID %q, got %q", expectedID, got.ID)
	}
	if got.Original != original {
		t.Fatalf("Shorten() expected Original %q, got %q", original, got.Original)
	}
}

func TestURLService_Shorten_SaveError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockURLReaderWriter(ctrl)
	svc := NewURLService(mockRepo)

	original := "https://example.com/error"
	saveErr := errors.New("save failed")

	mockRepo.
		EXPECT().
		Save(gomock.Any()).
		Return(saveErr).
		Times(1)

	got, err := svc.Shorten(original)
	if err == nil {
		t.Fatalf("Shorten() expected error, got nil")
	}
	if got != nil {
		t.Fatalf("Shorten() expected nil url on error, got: %+v", got)
	}
}

func TestURLService_Retrieve_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockURLReaderWriter(ctrl)
	svc := NewURLService(mockRepo)

	id := "abc12345"
	expected := &URL{ID: id, Original: "https://example.com/x"}

	mockRepo.
		EXPECT().
		Get(id).
		Return(expected, nil).
		Times(1)

	got, err := svc.Retrieve(id)
	if err != nil {
		t.Fatalf("Retrieve() unexpected error: %v", err)
	}
	if got == nil {
		t.Fatalf("Retrieve() returned nil url")
	}
	if got.ID != expected.ID || got.Original != expected.Original {
		t.Fatalf("Retrieve() expected %+v, got %+v", expected, got)
	}
}

func TestURLService_Retrieve_GetError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockURLReaderWriter(ctrl)
	svc := NewURLService(mockRepo)

	id := "deadbeef"
	getErr := errors.New("not found")

	mockRepo.
		EXPECT().
		Get(id).
		Return(nil, getErr).
		Times(1)

	got, err := svc.Retrieve(id)
	if err == nil {
		t.Fatalf("Retrieve() expected error, got nil")
	}
	if got != nil {
		t.Fatalf("Retrieve() expected nil url on error, got: %+v", got)
	}
}
