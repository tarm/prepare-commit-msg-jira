/*
Documentation: http://godoc.org/github.com/tarm/prepare-commit-msg-jira

Travis Status: https://travis-ci.org/tarm/prepare-commit-msg-jira

Prepares a git commit message from a Jira issue by fetching the
summary and description and formatting in a way that is appropiate for
a git commit message.  It identifies the Jira ticket number from the
name of the git branch that you are on.

You should install it in your .git/hooks/ directory as 'prepare-commit-msg'.

It uses the following git configuration items:

 - jira.url - this is the base url of the jira service.  It should
              have a trailing '/'

 - jira.regexp - this is a regular expression that extracts the issue
                 name from the git branch name.

 - jira.jsessionid - this is the value of the JSESSIONID cookie that
                     jira sets after you login.  Get this from your
                     browser.

You should set these git configuration items like this:
 git config --add jira.url https://jira.atlassian.com

If it cannot extract the ticket number from the branch or it cannot
find the ticket number in Jira, then the default git commit message
will be used.

It only prepares a commit message for "normal" commits.  It leaves the
git default commit message for ammend, --C, merge, and -m"msg" style
commits.

It applies the following format for the git commit message:

    TICKET-12345: The summary line

    Descriptive text is wrapped at 76 characters.  This is a really long
    line to get the point across.

        Text that is within the {noformat} directive is included without wrapping but is indented by 4 spaces

Patches (with unit tests) are welcome to improve authentication,
enhance the formatting, etc.
*/
package main

//go:generate godoc2md -o README.md github.com/tarm/prepare-commit-msg-jira

/////////////////////////////////////////
// Fixup the readme to add the badges: //
/////////////////////////////////////////
//go:generate sed -i "s|Documentation: .*|[![GoDoc](https://godoc.org/github.com/tarm/prepare-commit-msg-jira?status.svg)](http://godoc.org/github.com/tarm/prepare-commit-msg-jira)|" README.md
//go:generate sed -i "s|Travis Status: .*|[![Build Status](https://travis-ci.org/tarm/prepare-commit-msg-jira.svg?branch=master)](https://travis-ci.org/tarm/prepare-commit-msg-jira)|" README.md
