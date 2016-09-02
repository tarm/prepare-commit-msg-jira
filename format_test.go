package main

import (
	"encoding/hex"
	"testing"

	jira "github.com/andygrunwald/go-jira"
)

const expFormat = `EXAM-1234: This is an example summary

Here is the description. It has some lines that are seperated by only a single \n. They will 
be combined together and wrapped. 

A double newline should be it's own paragraph 

    A noformat section should be intented

`

func TestFormat(t *testing.T) {
	issue := &jira.Issue{
		Key: "EXAM-1234",
		Fields: &jira.IssueFields{
			Summary: "This is an example summary",
			Description: `Here is the description.
It has some lines that are seperated by only a single \n.  They will be combined together and wrapped.

A double newline should be it's own paragraph

{noformat}
A noformat section should be intented
{noformat}
`,
		},
	}

	out := formatForGit(issue)

	if string(out) != expFormat {
		t.Logf("%s\n", out)
		t.Fatalf("Bad formatting.  Expected:\n\n%v\nGot\n\n%v", hex.Dump([]byte(expFormat)), hex.Dump(out))
	}
}

func TestFetch(t *testing.T) {
	issue, err := fetchIssue("JRA-808", "9D8B6D83F4E00F3A83EC0A76304DE343.node1", "https://jira.atlassian.com")
	if err != nil {
		t.Fatal(err)
	}
	if issue == nil {
		t.Fatal("Should not have nil issue")
	}
	if issue.Fields.Summary != "jsessionid trouble" {
		t.Fatalf("Unexpected summary %v", issue.Fields.Summary)
	}
	if len(issue.Fields.Description) != 1230 {
		t.Fatalf("Got bad description length %v", len(issue.Fields.Description))
	}
}
