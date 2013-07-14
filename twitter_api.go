package twitter_api

import (
  "github.com/NOX73/go-oauth"
  "net/http"
  "bufio"
  "bytes"
)

type Credentials struct {
  oauth_consumer_key string
  oauth_token string
  oauth_consumer_secret string
  oauth_token_secret string
}

const (
  StatusMessageKind = 0
  TweetMessageKind = 1
)

type Messager interface {
  Kind() int
  Tweet() *TweetMessage
  Status() *StatusMessage
}

type StatusMessage struct {
  StatusCode int
  Error bool
}

func (m *StatusMessage) Kind() int {
  return StatusMessageKind
}

func (m *StatusMessage) Status() *StatusMessage {
  return m
}

func (m *StatusMessage) Tweet() *TweetMessage {
  return &TweetMessage{}
}



type TweetMessage struct {
  Body string
}

func (m *TweetMessage) Kind() int {
  return TweetMessageKind
}

func (m *TweetMessage) Tweet() *TweetMessage {
  return m
}

func (m *TweetMessage) Status() *StatusMessage {
  return &StatusMessage{} 
}

func TwitterStream (ch chan Messager, credentials *Credentials, params map[string]string){
  c := oauth.NewCredentials(credentials.oauth_consumer_key, credentials.oauth_token, credentials.oauth_consumer_secret, credentials.oauth_token_secret)

  r, _ := oauth.NewRequest("POST", "https://stream.twitter.com/1.1/statuses/filter.json", params, c)

  client := http.Client{}
  resp, _ := client.Do(r.HttpRequest())

  message := &StatusMessage{
    StatusCode: resp.StatusCode,
    Error: false,
  }

  ch <- message

  body_reader := bufio.NewReader(resp.Body)
  for {
    var part []byte //Part of line
    var prefix bool //Flag. Readln readed only part of line.
    var sep []byte // Separator for Join function

    parts := make([][]byte, 0, 5)

    part, prefix, err := body_reader.ReadLine()
    parts = append(parts, part)

    for !prefix {
      part, prefix, err = body_reader.ReadLine()
      parts = append(parts, part)
    }

    if err != nil { break }

    tweet := &TweetMessage{
      Body: string(bytes.Join(parts, sep)),
    }

    ch <- tweet
  }
}
