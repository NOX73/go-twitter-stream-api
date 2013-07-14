package twitter_api

import (
  "github.com/NOX73/go-oauth"
  "net/http"
  "bufio"
  "errors"
)

const (
  NewRequestMethod = "POST"
  NewRequestURL = "https://stream.twitter.com/1.1/statuses/filter.json"
  TwitterStreamApiConnestionError = "Error: Response status code not 200."
)

type Credentials struct {
  OauthConsumerKey string
  OauthToken string
  OauthConsumerSecret string
  OauthTokenSecret string
}

type Message struct {
  Error error
  Response *http.Response
  Tweet *Tweet
}

type Tweet struct {
  Body string
}

func NewCredentials(consumer_key, token, consumer_secret, token_secret string) *Credentials {
  return &Credentials{consumer_key, token, consumer_secret, token_secret}
}


func TwitterStream (ch chan Message, credentials *Credentials, params map[string]string){
  var message Message

  c := oauth.NewCredentials(credentials.OauthConsumerKey, credentials.OauthToken, credentials.OauthConsumerSecret, credentials.OauthTokenSecret)

  r, _ := oauth.NewRequest(NewRequestMethod, NewRequestURL, params, c)

  client := http.Client{}
  resp, err := client.Do(r.HttpRequest())

  if err == nil{
    err = CheckError(resp)
  }

  if err != nil {
    message = Message{
      Error:err,
      Response: resp,
    }

    ch <- message
    return
  }

  body_reader := bufio.NewReader(resp.Body)

  for {
    var part []byte //Part of line
    var prefix bool //Flag. Readln readed only part of line.

    part, prefix, err := body_reader.ReadLine()
    if err != nil { break }

    buffer := append([]byte(nil), part...)

    for prefix && err == nil {
      part, prefix, err = body_reader.ReadLine()
      buffer = append(buffer, part...)
    }
    if err != nil { break }

    tweet := &Tweet{
      Body: string(buffer),
    }

    message = Message{
      Response: resp,
      Tweet: tweet,
    }

    ch <- message
  }
}

func CheckError(r *http.Response) error {
  if r.StatusCode != 200{
    return errors.New(TwitterStreamApiConnestionError)
  }
  return nil
}
