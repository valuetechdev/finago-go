package payday

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Token struct {
	AccessToken string
	ExpiresAt   time.Time
}

func (c *PaydayClient) Authenticate() error {
	url := fmt.Sprintf("https://payroll.24sevenoffice.com/api/auth?token=%s", c.secret)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	type tokenRes struct {
		AccessToken string `json:"accessToken"`
		ExpiresOn   string `json:"expiresOn"`
	}
	payrollTokenRes := tokenRes{}
	if err = json.Unmarshal(body, &payrollTokenRes); err != nil {
		return err
	}

	expiresOnSplitted := strings.Split(payrollTokenRes.ExpiresOn, "+")
	expiresAt, err := time.Parse("2006-01-02T15:04:05", expiresOnSplitted[0])
	if err != nil {
		return err
	}

	c.token = &Token{
		AccessToken: payrollTokenRes.AccessToken,
		ExpiresAt:   expiresAt,
	}

	return nil
}

// Returns token
func (c *PaydayClient) GetToken() *Token {
	return c.token
}

// Sets token
func (c *PaydayClient) SetToken(token *Token) {
	c.token = token
}

// Returns wether the token is valid or not
func (c *PaydayClient) IsTokenValid() bool {
	if c.token == nil {
		return false
	}
	return time.Now().Before(c.token.ExpiresAt)
}

// Intercept the req with an "Authorization: bearer <accessToken>" header.
func (c *PaydayClient) Intercept(ctx context.Context, req *http.Request) error {
	req.Header.Set("Authorization", "Bearer "+c.token.AccessToken)

	return nil
}
