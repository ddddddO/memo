package handlers

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/mock"
)

func TestAuthHandler_Login(t *testing.T) {
	sessionSec := "sessionsecret"
	store := sessions.NewCookieStore([]byte(sessionSec))

	tests := []struct {
		name       string
		server     func() *httptest.Server
		body       func() io.Reader
		wantResp   string
		wantStatus int
	}{
		{
			name: "success",
			server: func() *httptest.Server {
				m := mockUserUseCase{}
				m.On("FetchUserID", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(7, nil)
				authHandler := NewAuthHandler(m)
				return httptest.NewServer(authHandler.Login(store))
			},
			body: func() io.Reader {
				v := url.Values{}
				v.Add("name", "testname")
				v.Add("passwd", "testpasswd")
				return strings.NewReader(v.Encode())
			},
			wantResp:   `{"user_id":7}`,
			wantStatus: http.StatusOK,
		},
		{
			name: "only name param",
			server: func() *httptest.Server {
				m := mockUserUseCase{}
				authHandler := NewAuthHandler(m)
				return httptest.NewServer(authHandler.Login(store))
			},
			body: func() io.Reader {
				v := url.Values{}
				v.Add("name", "testname")
				return strings.NewReader(v.Encode())
			},
			wantResp:   `{"message":"empty key 'passwd'"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "only passwd param",
			server: func() *httptest.Server {
				m := mockUserUseCase{}
				authHandler := NewAuthHandler(m)
				return httptest.NewServer(authHandler.Login(store))
			},
			body: func() io.Reader {
				v := url.Values{}
				v.Add("passwd", "testpasswd")
				return strings.NewReader(v.Encode())
			},
			wantResp:   `{"message":"empty key 'name'"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "UseCase failed",
			server: func() *httptest.Server {
				m := mockUserUseCase{}
				m.On("FetchUserID", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(0, errors.New("error !"))
				authHandler := NewAuthHandler(m)
				return httptest.NewServer(authHandler.Login(store))
			},
			body: func() io.Reader {
				v := url.Values{}
				v.Add("name", "testfailname")
				v.Add("passwd", "testfailpasswd")
				return strings.NewReader(v.Encode())
			},
			wantResp:   `{"message":"failed"}`,
			wantStatus: http.StatusInternalServerError,
		},
	}

	client := &http.Client{}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ts := tt.server()
			req, err := http.NewRequest(http.MethodPost, ts.URL+"/auth", tt.body())
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			resp, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			rBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}

			if resp.StatusCode != tt.wantStatus {
				t.Error("unexpected status code")
			}
			if string(rBody) != tt.wantResp {
				t.Error("unexpected response")
			}
		})
	}
}

type mockUserUseCase struct {
	mock.Mock
}

func (m mockUserUseCase) FetchUserID(name, passwd string) (int, error) {
	ret := m.Called(name, passwd)
	return ret.Get(0).(int), ret.Error(1)
}
