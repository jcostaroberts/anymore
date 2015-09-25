package main

import (
	"fmt"
	"github.com/kurrik/oauth1a"
	"github.com/kurrik/twittergo"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func check(err error, problem string) {
	if err == nil {
		return
	}
	fmt.Printf(problem, ": %v\n", err)
	os.Exit(1)
}

func getTwitterClient() (client *twittergo.Client, err error) {
	credentials, err := ioutil.ReadFile("CREDENTIALS")
	check(err, "Could not read CREDENTIALS")

	lines := strings.Split(string(credentials), "\n")
	config := &oauth1a.ClientConfig{
		ConsumerKey:    lines[0],
		ConsumerSecret: lines[1],
	}
	user := oauth1a.NewAuthorizedConfig(lines[2], lines[3])
	client = twittergo.NewClient(config, user)
	return client, err
}

func displayTweets(results *twittergo.SearchResults) {
	for i, tweet := range results.Statuses() {
		user := tweet.User()
		fmt.Printf("%v.) %v\n", i+1, tweet.Text())
		fmt.Printf("From %v (@%v) ", user.Name(), user.ScreenName())
		fmt.Printf("at %v\n\n", tweet.CreatedAt().Format(time.RFC1123))
	}
}

func displayRateInfo(resp *twittergo.APIResponse) {
	if resp.HasRateLimit() {
		fmt.Printf("Rate limit:           %v\n", resp.RateLimit())
		fmt.Printf("Rate limit remaining: %v\n", resp.RateLimitRemaining())
		fmt.Printf("Rate limit reset:     %v\n", resp.RateLimitReset())
	} else {
		fmt.Printf("Could not parse rate limit from response.\n")
	}
}

func main() {
	var (
		err     error
		client  *twittergo.Client
		req     *http.Request
		resp    *twittergo.APIResponse
		results *twittergo.SearchResults
	)

	// Use Twitter API credentials to start up client
	client, err = getTwitterClient()
	check(err, "Could not parse CREDENTIALS file")

	// Build request
	query := url.Values{}
	query.Set("q", "anymore")
	//query.Set("q", "anymore%20-\"not\"%20-\"can%27t\"%20-\"won%27t\"%20-\"isn%27t\"%20-\"don%27t\"%20-\"%3F\"")
	query.Set("count", "100")
	u := fmt.Sprintf("/1.1/search/tweets.json?%v", query.Encode())
	req, err = http.NewRequest("GET", u, nil)
	check(err, "Could not parse request")

	// Send request
	resp, err = client.SendRequest(req)
	check(err, "Could not send request")

	// Parse response
	results = &twittergo.SearchResults{}
	err = resp.Parse(results)
	check(err, "Could not parse response")

	// Display tweets
	displayTweets(results)
	displayRateInfo(resp)
}
