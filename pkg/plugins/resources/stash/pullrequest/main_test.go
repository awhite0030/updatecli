package pullrequest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	stashclient "github.com/updatecli/updatecli/pkg/plugins/resources/stash/client"
	stashscm "github.com/updatecli/updatecli/pkg/plugins/scms/stash"
)

func TestNew(t *testing.T) {
	testData := []struct {
		name                      string
		spec                      Spec
		scm                       *stashscm.Stash
		expectedStashOwner        string
		expectedStashRepository   string
		expectedStashSourceBranch string
		expectedStashTargetBranch string
		wantErr                   bool
	}{
		{
			name: "Test parameter inheritance 2",
			spec: Spec{},
			scm: &stashscm.Stash{
				Spec: stashscm.Spec{
					Spec: stashclient.Spec{
						URL: "stash.updatecli.io",
					},
					Repository: "updatecli-test",
					Branch:     "v2",
					Owner:      "tartempion",
				},
			},
			expectedStashOwner:        "tartempion",
			expectedStashRepository:   "updatecli-test",
			expectedStashSourceBranch: "v2",
			expectedStashTargetBranch: "v2",
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			g, gotErr := New(tt.spec, tt.scm)

			if tt.wantErr {
				require.Error(t, gotErr)
			} else {
				require.NoError(t, gotErr)
			}

			assert.Equal(t, g.Owner, tt.expectedStashOwner)
			assert.Equal(t, g.Repository, tt.expectedStashRepository)
			assert.Equal(t, g.SourceBranch, tt.expectedStashSourceBranch)
			assert.Equal(t, g.TargetBranch, tt.expectedStashTargetBranch)
		})
	}
}
