package twitter_api

import (
  . "launchpad.net/gocheck"
  "testing"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }
type OAuthSuite struct{}
var _ = Suite(&OAuthSuite{})

func (s *OAuthSuite) TestCreateClient(c *C) {

  ch := make(chan Message)
  params := make(map[string]string, 1)

  credentials := NewCredentials("XjY7q0CYwRxSBzCpUeRDzQ", "214373359-jn77FNlrKEajR4Gpp9l5msb1KXCGXZ7QeJPtt5TF", "cuseCPmxY4taUEmouOhXIvR7MVSUWdRKjKHvHKgVvOk", "tO5hW1ye3myBnT78DspVbTKWFgadvKeU1EOiV3o5Tg")

  params["track"] = "twitter"

  go TwitterStream(ch, credentials, params)

  message := <- ch

  c.Assert(message.Response.StatusCode, Equals, 200)
  c.Assert(message.Tweet, NotNil)
  c.Assert(message.Tweet.Body, NotNil)
  c.Assert(message.Tweet.Text(), NotNil)
  c.Assert(message.Tweet.UserID(), NotNil)
  c.Assert(message.Tweet.UserName(), NotNil)
}
