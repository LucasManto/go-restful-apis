package handlers

import (
	"net/http"
	"net/url"
	"os"
	"strconv"
	"testing"

	"github.com/LucasManto/go-restful-apis/user"
	"gopkg.in/mgo.v2/bson"
)

const (
	dbPath = "users.db"
)

type response struct {
	header http.Header
	code   int
	body   []byte
}

type mockWriter response

func newMockWriter() *mockWriter {
	return &mockWriter{
		body:   []byte{},
		header: http.Header{},
	}
}

func (mw *mockWriter) Write(b []byte) (int, error) {
	mw.body = make([]byte, len(b))
	copy(mw.body, b)
	return len(b), nil
}

func (mw *mockWriter) WriteHeader(code int) {
	mw.code = code
}

func (mw *mockWriter) Header() http.Header {
	return mw.header
}

func TestMain(m *testing.M) {
	m.Run()
	os.Remove(dbPath)
}

func prepDb(n int) error {
	os.Remove(dbPath)
	for i := 0; i < n; i++ {
		u := &user.User{
			ID:   bson.NewObjectId(),
			Name: "Lucas_" + strconv.Itoa(i),
			Role: "Software Engineer",
		}
		err := u.Save()
		if err != nil {
			return err
		}
	}
	return nil
}

func makeRequest() (*http.Request, error) {
	u, err := url.Parse("/users")
	if err != nil {
		return nil, err
	}
	req := &http.Request{
		URL:    u,
		Header: http.Header{},
		Method: http.MethodGet,
	}
	return req, nil
}

func getAll(b *testing.B, r *http.Request) {
	prepDb(100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		mw := newMockWriter()
		b.StartTimer()
		UsersRouter(mw, r)
	}
}

func BenchmarkGetAllNonCached(b *testing.B) {
	r, err := makeRequest()
	if err != nil {
		b.Fatal(err)
	}
	r.Header.Add("Cache-Control", "no-cache")
	getAll(b, r)
}

func BenchmarkGetAllCached(b *testing.B) {
	r, err := makeRequest()
	if err != nil {
		b.Fatal(err)
	}
	getAll(b, r)
}
