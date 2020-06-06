package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"io"
	"io/ioutil"
	"testing"
	"strings"

	"github.com/gorilla/sessions"
)

func TestAuthHandler_Login(t *testing.T) {
	sessionSec := "sessionsecret"
	store := sessions.NewCookieStore([]byte(sessionSec))

	mock := mockUserUseCase{}
	authHandler := NewAuthHandler(mock)
	server := httptest.NewServer(authHandler.Login(store))
	defer server.Close()

	client := &http.Client{}

	// sucsess
	v := url.Values{}
	v.Add("name", "testname")
	v.Add("passwd", "testpasswd")
	body := strings.NewReader(v.Encode())

	// only name param
	v2 := url.Values{}
	v2.Add("name", "testname")
	body2 := strings.NewReader(v2.Encode())

	// only passwd param
	v3 := url.Values{}
	v3.Add("passwd", "testpasswd")
	body3 := strings.NewReader(v3.Encode())

	// UseCase failed
	mockFail := mockFailUserUseCase{}
	authFailHandler := NewAuthHandler(mockFail)
	serverFail := httptest.NewServer(authFailHandler.Login(store))
	defer serverFail.Close()
	body4 := strings.NewReader(v.Encode())
	
	tests := []struct{
		name       string
		server     *httptest.Server
		body       io.Reader
		wantResp   string
		wantStatus int
	}{
		{
			name: "success",
			server: server,
			body: body,
			wantResp: `{"user_id":7}`,
			wantStatus: http.StatusOK,
		},
		{
			name: "only name param",
			server: server,
			body: body2,
			wantResp: `{"message":"empty key 'passwd'"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "only passwd param",
			server: server,
			body: body3,
			wantResp: `{"message":"empty key 'name'"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "UseCase failed",
			server: serverFail,
			body: body4,
			wantResp: `{"message":"failed"}`,
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Log(tt.name)
		req, err := http.NewRequest(http.MethodPost, tt.server.URL + "/auth", tt.body)
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
			t.Log(resp.StatusCode)
		}
		if string(rBody) != tt.wantResp {
			t.Error("unexpected response")
			t.Log(string(rBody))
		}
	
	}
}

type mockUserUseCase struct{}

func (m mockUserUseCase) FetchUserID(name, passwd string) (int, error) {
	return 7, nil
}

type mockFailUserUseCase struct{}

func (mf mockFailUserUseCase) FetchUserID(name, passwd string) (int, error) {
	return 0, errors.New("fail test")
}