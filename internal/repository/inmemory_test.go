package repository

import (
	"testing"
)

type item struct {
	ID    string
	Value int
}

func keyFn(i *item) string { return i.ID }

func TestInMemoryRepository_Get(t *testing.T) {
	repo := NewInMemoryRepository(keyFn)

	// Preload data for "found" scenario
	existing := &item{ID: "a1", Value: 42}
	if err := repo.Save(existing); err != nil {
		t.Fatalf("Save(existing) error = %v", err)
	}

	tests := []struct {
		name      string
		id        string
		wantErr   bool
		errMsg    string
		wantNil   bool
		wantValue int
		wantSame  *item
	}{
		{
			name:      "found_should_return_value_and_do_not_return_error",
			id:        "a1",
			wantErr:   false,
			wantNil:   false,
			wantValue: 42,
			wantSame:  existing,
		},
		{
			name:    "not_found_should_return_error_and_nil_value",
			id:      "missing",
			wantErr: true,
			errMsg:  "not found",
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.Get(tt.id)
			if tt.wantErr && err == nil {
				t.Fatalf("Get() error = nil, want error")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("Get() error = %v, want nil", err)
			}
			if tt.wantErr && err != nil && tt.errMsg != "" && err.Error() != tt.errMsg {
				t.Fatalf("Get() error = %q, want %q", err.Error(), tt.errMsg)
			}
			if tt.wantNil && got != nil {
				t.Fatalf("Get() = %v, want nil", got)
			}
			if !tt.wantNil && got == nil {
				t.Fatalf("Get() = nil, want non-nil")
			}
			if got != nil {
				if tt.wantValue != 0 && got.Value != tt.wantValue {
					t.Fatalf("Get() Value = %d, want %d", got.Value, tt.wantValue)
				}
				if tt.wantSame != nil && got != tt.wantSame {
					t.Fatalf("Get() returned pointer %p, want %p", got, tt.wantSame)
				}
			}
		})
	}
}

func TestInMemoryRepository_Delete(t *testing.T) {
	tests := []struct {
		name           string
		preSave        *item
		deleteID       string
		wantDelErr     bool
		postGetID      string
		wantPostGetErr bool
		postGetErrMsg  string
		wantPostGetNil bool
	}{
		{
			name:           "delete_existing_then_get_should_return_error_and_nil_value",
			preSave:        &item{ID: "x1", Value: 1},
			deleteID:       "x1",
			wantDelErr:     false,
			postGetID:      "x1",
			wantPostGetErr: true,
			postGetErrMsg:  "not found",
			wantPostGetNil: true,
		},
		{
			name:           "delete_non_existing_is_no_op",
			preSave:        nil,
			deleteID:       "nope",
			wantDelErr:     false,
			postGetID:      "nope",
			wantPostGetErr: true,
			postGetErrMsg:  "not found",
			wantPostGetNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewInMemoryRepository(keyFn)
			if tt.preSave != nil {
				if err := repo.Save(tt.preSave); err != nil {
					t.Fatalf("Save(preSave) error = %v", err)
				}
			}

			if err := repo.Delete(tt.deleteID); (err != nil) != tt.wantDelErr {
				t.Fatalf("Delete() error = %v, wantDelErr=%v", err, tt.wantDelErr)
			}

			got, err := repo.Get(tt.postGetID)
			if tt.wantPostGetErr && err == nil {
				t.Fatalf("Get() after Delete error = nil, want error")
			}
			if !tt.wantPostGetErr && err != nil {
				t.Fatalf("Get() after Delete error = %v, want nil", err)
			}
			if tt.wantPostGetErr && err != nil && tt.postGetErrMsg != "" && err.Error() != tt.postGetErrMsg {
				t.Fatalf("Get() after Delete error = %q, want %q", err.Error(), tt.postGetErrMsg)
			}
			if tt.wantPostGetNil && got != nil {
				t.Fatalf("Get() after Delete = %v, want nil", got)
			}
		})
	}
}

func TestInMemoryRepository_Save_Upsert(t *testing.T) {
	repo := NewInMemoryRepository(keyFn)

	first := &item{ID: "k1", Value: 10}
	if err := repo.Save(first); err != nil {
		t.Fatalf("Save(first) error = %v", err)
	}

	second := &item{ID: "k1", Value: 20}
	if err := repo.Save(second); err != nil {
		t.Fatalf("Save(second) error = %v", err)
	}

	got, err := repo.Get("k1")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if got == nil {
		t.Fatalf("Get() = nil, want non-nil")
	}
	if got.Value != 20 {
		t.Fatalf("Get() Value = %d, want %d", got.Value, 20)
	}
	if got != second {
		t.Fatalf("Get() returned pointer %p, want %p (latest saved)", got, second)
	}
}
