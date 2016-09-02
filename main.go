package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strings"

	jira "github.com/andygrunwald/go-jira"
)

func main() {
	// Don't modify if a pre-existing message exists
	if len(os.Args) > 2 {
		return
	}

	jsessionid, err := getGitConfig("jira.jsessionid")
	if err != nil {
		log.Fatal("Please set jira.jsessionid in you git config")
	}
	jurl, err := getGitConfig("jira.url")
	if err != nil {
		log.Fatal("Please set jira.url in your git config")
	}
	jregex, err := getGitConfig("jira.regexp")
	if err != nil {
		log.Fatal("Please set jira.regexp in your git config")
	}
	re, err := regexp.Compile(jregex)
	if err != nil {
		log.Fatal(err)
	}

	// Get the name from the branch
	out, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").CombinedOutput()
	if err != nil {
		log.Fatal("Could not get branch name", err)
	}
	name := string(re.Find(out))
	if name == "" {
		// We still want to allow regular committing
		fmt.Println("No matching branch found")
		return
	}

	issue, err := fetchIssue(name, jsessionid, jurl)
	if err != nil {
		// We still want to allow regular committing
		fmt.Println(err)
		return
	}

	err = ioutil.WriteFile(os.Args[1], formatForGit(issue), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func getGitConfig(key string) (string, error) {
	out, err := exec.Command("git", "config", "--get", key).CombinedOutput()
	return string(bytes.TrimSpace(out)), err
}

const lineLen = 76

// 1) Puts the summary at the top
// 2) wraps long lines at 76 characters
// 3) prettifies some jira formatting
//      - {noformat} gets indented by 4 spaces
//      - TODO more things like tables, quotes, and numbered lists
func formatForGit(i *jira.Issue) []byte {
	var buf bytes.Buffer

	// Summary
	fmt.Fprintf(&buf, "%v: %v\n\n", i.Key, i.Fields.Summary)

	// Wrap some lines
	str := strings.Replace(i.Fields.Description, "\r", "", -1)
	pars := strings.Split(str, "\n")
	count := 0
	noformat := false
	for i := range pars {
		if noformat {
			if pars[i] == "{noformat}" {
				noformat = false
			} else {
				fmt.Fprintf(&buf, "    %v\n", pars[i])
			}
			continue
		}
		if pars[i] == "{noformat}" {
			if count != 0 {
				buf.WriteString("\n\n")
				count = 0
			}
			noformat = true
			continue
		}

		if len(pars[i]) == 0 {
			// Double newline
			if count != 0 {
				buf.WriteString("\n")
				count = 0
			}
			buf.WriteString("\n")
			continue
		}
		firstChar := pars[i][0]
		switch {
		case firstChar >= 'a' && firstChar <= 'z':
		case firstChar >= 'A' && firstChar <= 'Z':
		case firstChar >= '0' && firstChar <= '9':
		default:
			if count != 0 {
				buf.WriteString("\n")
			}
			buf.WriteString(pars[i])
			buf.WriteString("\n")
			count = 0
			continue
		}
		words := strings.Fields(pars[i])
		for _, word := range words {
			if count+len(word) > lineLen {
				buf.WriteString("\n")
				count = 0
			}
			fmt.Fprintf(&buf, "%v ", word)
			count += len(word)
		}
	}

	return buf.Bytes()
}

func fetchIssue(name, jsessionid, jurl string) (*jira.Issue, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	jurlp, err := url.Parse(jurl)
	if err != nil {
		return nil, err
	}
	jar.SetCookies(jurlp, []*http.Cookie{{
		Name:  "JSESSIONID",
		Value: jsessionid,
	}})
	cl := &http.Client{
		Jar: jar,
	}

	c, err := jira.NewClient(cl, jurl)
	if err != nil {
		return nil, err
	}
	issue, resp, err := c.Issue.Get(name)
	if err != nil {
		if resp != nil {
			body, _ := ioutil.ReadAll(resp.Body)
			err = fmt.Errorf("%v %v", err, string(body))
		}
		return nil, err
	}

	return issue, nil
}
