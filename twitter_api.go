package twitter_api

import (
  "github.com/NOX73/go-oauth"
  "net/http"
  "fmt"
  "bufio"
  //"bytes"
  //"io"
)

type Tweet struct {
  id int
  text string
}

type Credentials struct {
  oauth_consumer_key string
  oauth_token string
  oauth_consumer_secret string
  oauth_token_secret string
}

func TwitterStream (ch chan Tweet, credentials *Credentials, params map[string]string){
  c := oauth.NewCredentials(credentials.oauth_consumer_key, credentials.oauth_token, credentials.oauth_consumer_secret, credentials.oauth_token_secret)

  r, _ := oauth.NewRequest("POST", "https://stream.twitter.com/1.1/statuses/filter.json?track=twitter", params, c)

  client := http.Client{}
  resp, _ := client.Do(r.HttpRequest())

  fmt.Println(resp.Status)

  body_reader := bufio.NewReader(resp.Body)
  for {
    line, prefix, _ := body_reader.ReadLine()
    fmt.Print(string(line))
    if !prefix {
      fmt.Println("\n----------------------------")
    }
  }
}
