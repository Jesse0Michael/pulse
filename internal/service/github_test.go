package service

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/google/go-github/v54/github"
	"github.com/jesse0michael/testhelpers/pkg/testserver"
)

func Test_eventActivity(t *testing.T) {
	repository := "test-repository"
	pushEvent := "PushEvent"
	pullRequestEvent := "PullRequestEvent"
	otherEvent := "Other"
	otherPayload := json.RawMessage(`{}`)
	pullRequest, _ := os.ReadFile("testdata/github_pullrequestevent_opened.json")
	pullRequestPayload := json.RawMessage(pullRequest)
	push, _ := os.ReadFile("testdata/github_pushevent.json")
	pushPayload := json.RawMessage(push)
	tests := []struct {
		name  string
		event *github.Event
		want  string
	}{
		{
			name:  "nil event",
			event: nil,
			want:  "",
		},
		{
			name:  "empty event",
			event: &github.Event{},
			want:  "",
		},
		{
			name: "other event",
			event: &github.Event{
				Type: &otherEvent,
				Repo: &github.Repository{
					Name: &repository,
				},
				RawPayload: &otherPayload,
			},
			want: "",
		},
		{
			name: "push event",
			event: &github.Event{
				Type: &pushEvent,
				Repo: &github.Repository{
					Name: &repository,
				},
				RawPayload: &pushPayload,
			},
			want: `pushed commits to repository test-repository
commit messages: feat: BREAKING CHANGE upgrade rest assured to v4
feat: export Serve method to start http listener

Remove client context and error channel.
Rely on the caller of the package to appropriately call Serve
chore: upgrade to go 1.21
feat: move to log/slog for logging
fix: use google/uuid package
test: Serve rest assured client in tests
chore: update license
feat: add NewClientServe function

that will create and serve the assured client
build: add docker labels
ci: update go version`,
		},
		{
			name: "pull request Event event",
			event: &github.Event{
				Type: &pullRequestEvent,
				Repo: &github.Repository{
					Name: &repository,
				},
				RawPayload: &pullRequestPayload,
			},
			want: `opened pull request in repository test-repository
body: [feat: BREAKING CHANGE upgrade rest assured to v4](https://github.com/Jesse0Michael/go-rest-assured/commit/b1d6cc11e692952793fa6484e540e46fbb1df11a)
[feat: export Serve method to start http listener](https://github.com/Jesse0Michael/go-rest-assured/commit/f51f3a59951b2316282ef1f698ef40abd0742d6b)
Remove client context and error channel.
Rely on the caller of the package to appropriately call Serve
[chore: upgrade to go 1.21](https://github.com/Jesse0Michael/go-rest-assured/commit/e90797d9c77f661fb0e64d440bc17ff46d872b1e)
[feat: move to log/slog for logging](https://github.com/Jesse0Michael/go-rest-assured/commit/95e02aca5a8ce5a065d920e56973fb606ef62fbf)
[fix: use google/uuid package](https://github.com/Jesse0Michael/go-rest-assured/commit/b7c3d2dc4939665ab817df2419a28db4376abb08)
[test: Serve rest assured client in tests](https://github.com/Jesse0Michael/go-rest-assured/commit/1e3cf4c77878750d09cfb555f1e850b39d32c39a)
[chore: update license](https://github.com/Jesse0Michael/go-rest-assured/commit/9fe24e3f833ed6149c3bca072bd837ecd6934dde)
[feat: add NewClientServe function](https://github.com/Jesse0Michael/go-rest-assured/commit/556a0d8c1de5d351624f9f671d5841dfca51b1f3)
[build: add docker labels](https://github.com/Jesse0Michael/go-rest-assured/commit/c5dbe4cd14961de9fba0dccd6717cf7c8deda84d)


update the go rest assured package to be more flexible.
update dependencies and documentation
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := eventActivity(tt.event); normalizeLineEndings(got) != tt.want {
				t.Errorf("eventActivity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parsePushEvent(t *testing.T) {
	b, _ := os.ReadFile("testdata/github_pushevent.json")
	var pushEvent github.PushEvent
	_ = json.Unmarshal(b, &pushEvent)
	tests := []struct {
		name  string
		event github.PushEvent
		want  string
	}{
		{
			name:  "empty event",
			event: github.PushEvent{},
			want:  "pushed commits to repository test-repository",
		},
		{
			name:  "valid event",
			event: pushEvent,
			want: `pushed commits to repository test-repository
commit messages: feat: BREAKING CHANGE upgrade rest assured to v4
feat: export Serve method to start http listener

Remove client context and error channel.
Rely on the caller of the package to appropriately call Serve
chore: upgrade to go 1.21
feat: move to log/slog for logging
fix: use google/uuid package
test: Serve rest assured client in tests
chore: update license
feat: add NewClientServe function

that will create and serve the assured client
build: add docker labels
ci: update go version`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parsePushEvent(tt.event, "test-repository"); got != tt.want {
				t.Errorf("parsePushEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parsePullRequestEvent(t *testing.T) {
	action := "opened"
	b, _ := os.ReadFile("testdata/github_pullrequestevent_opened.json")
	var opened github.PullRequestEvent
	_ = json.Unmarshal(b, &opened)
	b, _ = os.ReadFile("testdata/github_pullrequestevent_closed.json")
	var closed github.PullRequestEvent
	_ = json.Unmarshal(b, &closed)
	tests := []struct {
		name  string
		event github.PullRequestEvent
		want  string
	}{
		{
			name:  "empty event",
			event: github.PullRequestEvent{},
			want:  "",
		},
		{
			name:  "valid event: opened",
			event: opened,
			want: `opened pull request in repository test-repository
body: [feat: BREAKING CHANGE upgrade rest assured to v4](https://github.com/Jesse0Michael/go-rest-assured/commit/b1d6cc11e692952793fa6484e540e46fbb1df11a)
[feat: export Serve method to start http listener](https://github.com/Jesse0Michael/go-rest-assured/commit/f51f3a59951b2316282ef1f698ef40abd0742d6b)
Remove client context and error channel.
Rely on the caller of the package to appropriately call Serve
[chore: upgrade to go 1.21](https://github.com/Jesse0Michael/go-rest-assured/commit/e90797d9c77f661fb0e64d440bc17ff46d872b1e)
[feat: move to log/slog for logging](https://github.com/Jesse0Michael/go-rest-assured/commit/95e02aca5a8ce5a065d920e56973fb606ef62fbf)
[fix: use google/uuid package](https://github.com/Jesse0Michael/go-rest-assured/commit/b7c3d2dc4939665ab817df2419a28db4376abb08)
[test: Serve rest assured client in tests](https://github.com/Jesse0Michael/go-rest-assured/commit/1e3cf4c77878750d09cfb555f1e850b39d32c39a)
[chore: update license](https://github.com/Jesse0Michael/go-rest-assured/commit/9fe24e3f833ed6149c3bca072bd837ecd6934dde)
[feat: add NewClientServe function](https://github.com/Jesse0Michael/go-rest-assured/commit/556a0d8c1de5d351624f9f671d5841dfca51b1f3)
[build: add docker labels](https://github.com/Jesse0Michael/go-rest-assured/commit/c5dbe4cd14961de9fba0dccd6717cf7c8deda84d)


update the go rest assured package to be more flexible.
update dependencies and documentation
`,
		},
		{
			name: "empty event",
			event: github.PullRequestEvent{
				Action: &action,
			},
			want: "opened pull request in repository test-repository",
		},
		{
			name:  "valid event: closed",
			event: closed,
			want:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parsePullRequestEvent(tt.event, "test-repository"); normalizeLineEndings(got) != tt.want {
				t.Errorf("parsePullRequestEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGithub_UserActivity(t *testing.T) {
	listEvents, _ := os.ReadFile("testdata/github_listevents.json")
	tests := []struct {
		name    string
		token   string
		server  *testserver.Server
		want    string
		wantErr bool
	}{
		{
			name:  "list user activity",
			token: "test-token",
			server: testserver.New(
				testserver.Handler{Path: "/users/test-user/events", Status: http.StatusOK, Response: listEvents},
			),
			want: `pushed commits to repository Jesse0Michael/go-rest-assured
commit messages: ci: release from main
pushed commits to repository Jesse0Michael/go-rest-assured
commit messages: chore: change default branch to main
pushed commits to repository Jesse0Michael/go-rest-assured
commit messages: build: move docker labels
pushed commits to repository Jesse0Michael/go-rest-assured
commit messages: feat: BREAKING CHANGE upgrade rest assured to v4
feat: export Serve method to start http listener

Remove client context and error channel.
Rely on the caller of the package to appropriately call Serve
chore: upgrade to go 1.21
feat: move to log/slog for logging
fix: use google/uuid package
test: Serve rest assured client in tests
chore: update license
feat: add NewClientServe function

that will create and serve the assured client
build: add docker labels
ci: update go version
pushed commits to repository Jesse0Michael/go-rest-assured
commit messages: ci: update go version
opened pull request in repository Jesse0Michael/go-rest-assured
body: [feat: BREAKING CHANGE upgrade rest assured to v4](https://github.com/Jesse0Michael/go-rest-assured/commit/b1d6cc11e692952793fa6484e540e46fbb1df11a)
[feat: export Serve method to start http listener](https://github.com/Jesse0Michael/go-rest-assured/commit/f51f3a59951b2316282ef1f698ef40abd0742d6b)
Remove client context and error channel.
Rely on the caller of the package to appropriately call Serve
[chore: upgrade to go 1.21](https://github.com/Jesse0Michael/go-rest-assured/commit/e90797d9c77f661fb0e64d440bc17ff46d872b1e)
[feat: move to log/slog for logging](https://github.com/Jesse0Michael/go-rest-assured/commit/95e02aca5a8ce5a065d920e56973fb606ef62fbf)
[fix: use google/uuid package](https://github.com/Jesse0Michael/go-rest-assured/commit/b7c3d2dc4939665ab817df2419a28db4376abb08)
[test: Serve rest assured client in tests](https://github.com/Jesse0Michael/go-rest-assured/commit/1e3cf4c77878750d09cfb555f1e850b39d32c39a)
[chore: update license](https://github.com/Jesse0Michael/go-rest-assured/commit/9fe24e3f833ed6149c3bca072bd837ecd6934dde)
[feat: add NewClientServe function](https://github.com/Jesse0Michael/go-rest-assured/commit/556a0d8c1de5d351624f9f671d5841dfca51b1f3)
[build: add docker labels](https://github.com/Jesse0Michael/go-rest-assured/commit/c5dbe4cd14961de9fba0dccd6717cf7c8deda84d)


update the go rest assured package to be more flexible.
update dependencies and documentation

pushed commits to repository Jesse0Michael/fetcher
commit messages: fix: disable Instagram feed
pushed commits to repository Jesse0Michael/fetcher
commit messages: feat: add Untappd feed support
Merge branch 'main' of ssh://github.com/Jesse0Michael/fetcher
pushed commits to repository Jesse0Michael/fetcher
commit messages: chore: update twitter and insta package
`,
			wantErr: false,
		},
		{
			name: "failed list user activity",
			server: testserver.New(
				testserver.Handler{Path: "/users/test-user/events", Status: http.StatusNotFound, Response: []byte(`{"message": "Not Found"}`)},
			),
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGithub(GithubConfig{URL: tt.server.URL + "/", Token: tt.token})
			got, err := g.UserActivity(context.TODO(), "test-user")
			if (err != nil) != tt.wantErr {
				t.Errorf("Github.UserActivity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if normalizeLineEndings(got) != tt.want {
				t.Errorf("Github.UserActivity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func normalizeLineEndings(s string) string {
	// Replace all possible line ending sequences with '\n'
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	return s
}
