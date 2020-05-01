package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
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

	v := url.Values{}
	v.Add("name", "testname")
	v.Add("passwd", "testpasswd")
	body := strings.NewReader(v.Encode())
	req, err := http.NewRequest(http.MethodPost, server.URL + "/auth", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	
	tests := []struct{
		wantResp   string
		wantStatus int
		// TODO: body にパターン追加
	}{
		{
			wantResp: `{"user_id":7}`,
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
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
	
	}
}

type mockUserUseCase struct{}

func (m mockUserUseCase) FetchUserID(name, passwd string) (int, error) {
	return 7, nil
}