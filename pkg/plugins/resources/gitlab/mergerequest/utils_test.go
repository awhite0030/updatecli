package mergerequest

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	gitlabclient "github.com/updatecli/updatecli/pkg/plugins/resources/gitlab/client"
	gitlabscm "github.com/updatecli/updatecli/pkg/plugins/scms/gitlab"
)

func newMockServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/api/v4/projects/olblak/updatecli/repository/branches/main" {
			fmt.Fprintf(w, `{"name": "main"}`)
		} else if r.URL.Path == "/api/v4/projects/olblak/updatecli/merge_requests" {
			fmt.Fprintf(w, `[]`)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

func TestIsRemoteBranchExist(t *testing.T) {
	server := newMockServer()
	defer server.Close()

	testdata := []struct {
		name           string
		spec           Spec
		expectedResult bool
	}{
		{
			name: "Existing branch",
			spec: Spec{
				Owner:        "olblak",
				Repository:   "updatecli",
				SourceBranch: "main",
				TargetBranch: "main",
			},
			expectedResult: true,
		},
		{
			name: "Existing branch",
			spec: Spec{
				Owner:        "olblak",
				Repository:   "updatecli",
				SourceBranch: "donotexist",
				TargetBranch: "donotexist",
			},
			expectedResult: false,
		},
	}

	for _, td := range testdata {
		t.Run(td.name, func(t *testing.T) {
			scm := &gitlabscm.Gitlab{
				Spec: gitlabscm.Spec{
					Spec: gitlabclient.Spec{
						URL: server.URL,
					},
				},
			}
			gitlab, err := New(td.spec, scm)
			if err != nil {
				t.Fatalf("failed to create Gitlab instance: %v", err)
			}

			gotResult, err := gitlab.isRemoteBranchesExist()
			require.NoError(t, err)

			require.Equal(t, td.expectedResult, gotResult)

		})
	}
}
func TestFindExistingMR_Table(t *testing.T) {
	server := newMockServer()
	defer server.Close()

	testdata := []struct {
		name string
		spec Spec
	}{
		{
			name: "NoMergeRequest",
			spec: Spec{
				Owner:        "olblak",
				Repository:   "updatecli",
				SourceBranch: "donotexist",
				TargetBranch: "donotexist",
			},
		},
		{
			name: "NoOpenedMergeRequestBetweenMainBranches",
			spec: Spec{
				Owner:        "olblak",
				Repository:   "updatecli",
				SourceBranch: "main",
				TargetBranch: "main",
			},
		},
	}

	for _, td := range testdata {
		tc := td
		t.Run(tc.name, func(t *testing.T) {
			scm := &gitlabscm.Gitlab{
				Spec: gitlabscm.Spec{
					Spec: gitlabclient.Spec{
						URL: server.URL,
					},
				},
			}
			gitlab, err := New(tc.spec, scm)
			if err != nil {
				t.Fatalf("failed to create Gitlab instance: %v", err)
			}

			mr, err := gitlab.findExistingMR()
			require.NoError(t, err)
			require.Nil(t, mr)
		})
	}
}
