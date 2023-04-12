package urlsigner

import (
	"fmt"
	goalone "github.com/bwmarrin/go-alone"
	"strings"
	"time"
)

type Signer struct {
	Secret []byte
}

func (s *Signer) GenerateTokenFromString(token string) string {
	urlToSign := ""
	crypt := goalone.New(s.Secret, goalone.Timestamp)
	if strings.Contains(token, "?") {
		urlToSign = fmt.Sprintf("%s&hash=", token)
	} else {
		urlToSign = fmt.Sprintf("%s?hash=", token)
	}

	tokenBytes := crypt.Sign([]byte(urlToSign))
	t := string(tokenBytes)

	return t
}

func (s *Signer) VerifyToken(token string) bool {
	crypt := goalone.New(s.Secret, goalone.Timestamp)
	_, err := crypt.Unsign([]byte(token))
	
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func (s *Signer) IsExpired(token string, minutesUntilExpiry int) bool {
	crypt := goalone.New(s.Secret, goalone.Timestamp)
	ts := crypt.Parse([]byte(token))
	return time.Since(ts.Timestamp) > time.Duration(minutesUntilExpiry)*time.Minute
}
