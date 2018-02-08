// logview shows the CI logs for a particular build.
//
// The logs are queried from elasticsearch.
package logview

// TODO: error handling
// TODO: get rid of fatals

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"net/http"
	"reflect"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	elastic "gopkg.in/olivere/elastic.v5"

	"github.com/yourbase/yourbase/servicelocation/elasticloc"
)

// LogsIndices returns the ElasticSearch indices we consider relevant for logs
// search. Useful as an argument to elasticclient.Search.
func LogsIndices(client *elastic.Client) []string {
	// Logs are stored in the elasticsearch cluster. To find them, we
	// have to search using the "logstash-YYYY.MM.DD" indices, i.e: we need to
	// say explicitly which days we're searching for.
	// For now, we'll just search for the past 7 days because people rarely
	// look at build logs past this time.
	// Note that the name "logstash" is just historical because logstash is the
	// old school system that streamed log to ES. We're using fluent-bit/fluentd
	// instead.
	//
	// You can't simply add non-existent indices to the search query, apparently,
	// so we are going to filter them before asking.
	// Honestly I just coded this quickly, but there might be an elastic library
	// option for doing this automatically, or tolerating queries on non-existent
	// indices.
	// I assume this is efficient, if not we should test for the existence of
	// individual indices we want.
	indices, err := client.IndexNames()
	if err != nil {
		return nil
	}
	imap := map[string]bool{}
	for _, i := range indices {
		imap[i] = true
	}
	// Find indices for the last 7 days.
	days := []string{}
	today := time.Now()
	for i := 0; i < 7; i++ {
		day := today.Add(-time.Hour * time.Duration(24*i)).Format("logstash-2006.01.02")
		if imap[day] {
			days = append(days, day)
		}
	}
	return days
}

type Message struct {
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}

func viewLog(build string) (string, error) {
	client, err := elasticloc.Client()
	if err != nil {
		return "", fmt.Errorf("could not connect to elastic to search for logs for build %q: %v", build, err)
	}
	days := LogsIndices(client)

	termQuery := elastic.NewTermQuery("build", build)
	// TODO: Use scroll results instead of a large size.

	// Using @timestamp to sort because it's higher resolution than `time`.
	searchResult, err := client.Search(days...).Size(5000).Pretty(true).Sort("@timestamp", true).Query(termQuery).Do(context.Background())

	if err != nil {
		return "", fmt.Errorf("failed to start search on Elastic for build %q: %v", build, err)
	}
	// Number of hits
	if searchResult.Hits.TotalHits == 0 {
		// No hits
		log.Printf("Found no logs for %q", build)
		return "", nil
	}

	log.Printf("Found %d log entries\n", searchResult.Hits.TotalHits)

	var lines []Message
	var m Message
	for _, item := range searchResult.Each(reflect.TypeOf(m)) {
		m := item.(Message)
		lines = append(lines, m)
	}

	t := template.Must(template.New("").Parse(`<html><table>{{range .}}
		<tr><td><pre>{{.Time}} {{.Message}}</pre></td></tr>
		{{end}}`))

	var table bytes.Buffer
	if err := t.Execute(&table, lines); err != nil {
		return "", fmt.Errorf("Error executing log template: %v", err)
	}
	return table.String(), nil
}

// GitHubURL returns the URL that can be used to view the logs for a particular
// build.
func LogURL(owner, repo, branch, build string) string {
	// TODO: Parametrise the host for the appropriate YB universe.
	// TODO: SSL.
	return fmt.Sprintf("http://ci.microclusters.com/builds/github.com/%s/%s/%s/%s", owner, repo, branch, build)
}

func Handler(w http.ResponseWriter, req *http.Request) {
	// TODO: more robust method of naming path elements?

	parts := strings.Split(req.URL.Path, "/")
	if len(parts) != 7 {
		log.Printf("Logs request error: invalid path %#v", parts)
		http.NotFound(w, req)
		return
	}
	// /builds/github.com/yourbase/yourbase/add-yb/0b0096801e0c45e59737e43cc8286e84
	owner, repo, branch, build := parts[3], parts[4], parts[5], parts[6]
	log.Println("Searching for build logs from", owner, repo, branch, build)
	out, err := viewLog(build)
	if err != nil {
		http.Error(w, fmt.Sprintf("error searching for log: %v", err), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, out)
}
