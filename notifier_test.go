package notifier_test

import (
	"errors"
	"testing"

	"github.com/aceextensions/notifier"
)

// mockNotifier is a test double implementing the Notifier interface.
type mockNotifier struct {
	calledWith []struct{ subject, body string }
	returnErr  error
}

func (m *mockNotifier) Send(subject, body string) error {
	m.calledWith = append(m.calledWith, struct{ subject, body string }{subject, body})
	return m.returnErr
}

func TestMulti_SendsToAllNotifiers(t *testing.T) {
	a := &mockNotifier{}
	b := &mockNotifier{}

	multi := notifier.NewMulti(a, b)
	err := multi.Send("Test Subject", "Test Body")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(a.calledWith) != 1 {
		t.Errorf("expected notifier A to be called once, got %d", len(a.calledWith))
	}
	if len(b.calledWith) != 1 {
		t.Errorf("expected notifier B to be called once, got %d", len(b.calledWith))
	}
	if a.calledWith[0].subject != "Test Subject" {
		t.Errorf("expected subject 'Test Subject', got '%s'", a.calledWith[0].subject)
	}
}

func TestMulti_ContinuesOnError(t *testing.T) {
	a := &mockNotifier{returnErr: errors.New("a failed")}
	b := &mockNotifier{}

	multi := notifier.NewMulti(a, b)
	_ = multi.Send("Subject", "Body")

	// B must still be called even though A failed
	if len(b.calledWith) != 1 {
		t.Errorf("expected notifier B to be called even after A failed, got %d calls", len(b.calledWith))
	}
}

func TestMulti_ReturnsLastError(t *testing.T) {
	a := &mockNotifier{returnErr: errors.New("a failed")}
	b := &mockNotifier{returnErr: errors.New("b failed")}

	multi := notifier.NewMulti(a, b)
	err := multi.Send("Subject", "Body")

	if err == nil {
		t.Error("expected an error but got nil")
	}
}

func TestMulti_Empty(t *testing.T) {
	multi := notifier.NewMulti()
	err := multi.Send("Subject", "Body")
	if err != nil {
		t.Errorf("empty multi should not error, got: %v", err)
	}
}
