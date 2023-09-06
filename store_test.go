package sessions

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// Test for GH-8 for CookieStore
func TestGH8CookieStore(t *testing.T) {
	originalPath := "/"
	store := NewCookieStore()
	store.Options.Path = originalPath
	req, err := http.NewRequest("GET", "http://www.example.com", nil)
	if err != nil {
		t.Fatal("failed to create request", err)
	}

	session, err := store.New(req, "hello")
	if err != nil {
		t.Fatal("failed to create session", err)
	}

	store.Options.Path = "/foo"
	if session.Options.Path != originalPath {
		t.Fatalf("bad session path: got %q, want %q", session.Options.Path, originalPath)
	}
}

// 测试通过session str存取
func TestGH8FilesystemStoreByString(t *testing.T) {
	originalPath := "./__sessions"
	store := NewFilesystemStore(originalPath, []byte(os.Getenv("SESSION_KEY)")))
	store.Options.Path = originalPath
	store.Options.MaxAge = 86400

	// 存session
	t_sis := NewSession(store, "hero")
	t_sis.Values["user_id"] = 5
	t_sis.Options.MaxAge = 86400
	my_session, _ := store.SaveOnileSession(t_sis)

	fmt.Println("user_id: ", t_sis.Values["user_id"])

	// 取session
	kk, err := store.GetByToken(*my_session, "hero")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(kk.Values["user_id"])
	fmt.Println(kk)
}

// Test for GH-8 for FilesystemStore
func TestGH8FilesystemStore(t *testing.T) {
	originalPath := "/"
	store := NewFilesystemStore("")
	store.Options.Path = originalPath
	req, err := http.NewRequest("GET", "http://www.example.com", nil)
	if err != nil {
		t.Fatal("failed to create request", err)
	}

	session, err := store.New(req, "hello")
	if err != nil {
		t.Fatal("failed to create session", err)
	}

	store.Options.Path = "/foo"
	if session.Options.Path != originalPath {
		t.Fatalf("bad session path: got %q, want %q", session.Options.Path, originalPath)
	}
}

// Test for GH-2.
func TestGH2MaxLength(t *testing.T) {
	store := NewFilesystemStore("", []byte("some key"))
	req, err := http.NewRequest("GET", "http://www.example.com", nil)
	if err != nil {
		t.Fatal("failed to create request", err)
	}
	w := httptest.NewRecorder()

	session, err := store.New(req, "my session")
	if err != nil {
		t.Fatal("failed to create session", err)
	}

	session.Values["big"] = make([]byte, base64.StdEncoding.DecodedLen(4096*2))
	err = session.Save(req, w)
	if err == nil {
		t.Fatal("expected an error, got nil")
	}

	store.MaxLength(4096 * 3) // A bit more than the value size to account for encoding overhead.
	err = session.Save(req, w)
	if err != nil {
		t.Fatal("failed to Save:", err)
	}
}

// Test delete filesystem store with max-age: -1
func TestGH8FilesystemStoreDelete(t *testing.T) {
	store := NewFilesystemStore("", []byte("some key"))
	req, err := http.NewRequest("GET", "http://www.example.com", nil)
	if err != nil {
		t.Fatal("failed to create request", err)
	}
	w := httptest.NewRecorder()

	session, err := store.New(req, "hello")
	if err != nil {
		t.Fatal("failed to create session", err)
	}

	err = session.Save(req, w)
	if err != nil {
		t.Fatal("failed to save session", err)
	}

	session.Options.MaxAge = -1
	err = session.Save(req, w)
	if err != nil {
		t.Fatal("failed to delete session", err)
	}
}

// Test delete filesystem store with max-age: 0
func TestGH8FilesystemStoreDelete2(t *testing.T) {
	store := NewFilesystemStore("", []byte("some key"))
	req, err := http.NewRequest("GET", "http://www.example.com", nil)
	if err != nil {
		t.Fatal("failed to create request", err)
	}
	w := httptest.NewRecorder()

	session, err := store.New(req, "hello")
	if err != nil {
		t.Fatal("failed to create session", err)
	}

	err = session.Save(req, w)
	if err != nil {
		t.Fatal("failed to save session", err)
	}

	session.Options.MaxAge = 0
	err = session.Save(req, w)
	if err != nil {
		t.Fatal("failed to delete session", err)
	}
}
