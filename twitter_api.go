package twitter_api

import (
  "github.com/NOX73/go-oauth"
  "net/http"
  "bufio"
  "errors"
  "encoding/json"
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

// >> Tweet
// stream api message

type TweetJSON struct {
  Text string
  User struct {
    Id int
    Screen_name string
    Name string
    Description string
    Profile_image_url_https string
  }
}

type Tweet struct {
  Body string
  JSON *TweetJSON
}

func (t *Tweet) Text() string {
  if t.JSON == nil{ t.ParseJSON() }
  return t.JSON.Text
}

func (t *Tweet) UserID() int {
  if t.JSON == nil{ t.ParseJSON() }
  return t.JSON.User.Id
}

func (t *Tweet) UserName() string {
  if t.JSON == nil{ t.ParseJSON() }
  return t.JSON.User.Screen_name
}

func (t *Tweet) ParseJSON() {
  t.JSON = &TweetJSON{}
  _ = json.Unmarshal([]byte(t.Body), t.JSON) 
}

// << Tweet


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
