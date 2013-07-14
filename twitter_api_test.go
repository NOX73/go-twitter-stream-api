package twitter_api

import (
  . "launchpad.net/gocheck"
  "testing"
  "fmt"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }
type OAuthSuite struct{}
var _ = Suite(&OAuthSuite{})

func (s *OAuthSuite) TestCreateClient(c *C) {
  ch := make(chan Messager)
  params := make(map[string]string, 1)

  credentials := Credentials{
    oauth_consumer_key: "XjY7q0CYwRxSBzCpUeRDzQ",
    oauth_token: "214373359-jn77FNlrKEajR4Gpp9l5msb1KXCGXZ7QeJPtt5TF",
    oauth_consumer_secret: "cuseCPmxY4taUEmouOhXIvR7MVSUWdRKjKHvHKgVvOk",
    oauth_token_secret: "tO5hW1ye3myBnT78DspVbTKWFgadvKeU1EOiV3o5Tg",
  }

  params["track"] = "twitter"

  go TwitterStream(ch, &credentials, params)

  for {
    t := <- ch
    fmt.Println("Recieve message with kind = ", t.Kind())
    switch t.Kind() {
    case StatusMessageKind:
      fmt.Println("Status = ", t.Status().StatusCode)
    case TweetMessageKind:
      fmt.Println("Tweet body = ", t.Tweet().Body)
    }
  }
}
