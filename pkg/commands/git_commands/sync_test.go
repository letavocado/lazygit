package git_commands

import (
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/stretchr/testify/assert"
)

func TestSyncPush(t *testing.T) {
	type scenario struct {
		testName string
		opts     PushOpts
		test     func(oscommands.ICmdObj, error)
	}

	scenarios := []scenario{
		{
			testName: "Push with force disabled",
			opts:     PushOpts{Force: false},
			test: func(cmdObj oscommands.ICmdObj, err error) {
				assert.Equal(t, cmdObj.Args(), []string{"git", "push"})
				assert.NoError(t, err)
			},
		},
		{
			testName: "Push with force enabled",
			opts:     PushOpts{Force: true},
			test: func(cmdObj oscommands.ICmdObj, err error) {
				assert.Equal(t, cmdObj.Args(), []string{"git", "push", "--force-with-lease"})
				assert.NoError(t, err)
			},
		},
		{
			testName: "Push with force disabled, upstream supplied",
			opts: PushOpts{
				Force:          false,
				UpstreamRemote: "origin",
				UpstreamBranch: "master",
			},
			test: func(cmdObj oscommands.ICmdObj, err error) {
				assert.Equal(t, cmdObj.Args(), []string{"git", "push", "origin", "master"})
				assert.NoError(t, err)
			},
		},
		{
			testName: "Push with force disabled, setting upstream",
			opts: PushOpts{
				Force:          false,
				UpstreamRemote: "origin",
				UpstreamBranch: "master",
				SetUpstream:    true,
			},
			test: func(cmdObj oscommands.ICmdObj, err error) {
				assert.Equal(t, cmdObj.Args(), []string{"git", "push", "--set-upstream", "origin", "master"})
				assert.NoError(t, err)
			},
		},
		{
			testName: "Push with force enabled, setting upstream",
			opts: PushOpts{
				Force:          true,
				UpstreamRemote: "origin",
				UpstreamBranch: "master",
				SetUpstream:    true,
			},
			test: func(cmdObj oscommands.ICmdObj, err error) {
				assert.Equal(t, cmdObj.Args(), []string{"git", "push", "--force-with-lease", "--set-upstream", "origin", "master"})
				assert.NoError(t, err)
			},
		},
		{
			testName: "Push with remote branch but no origin",
			opts: PushOpts{
				Force:          true,
				UpstreamRemote: "",
				UpstreamBranch: "master",
				SetUpstream:    true,
			},
			test: func(cmdObj oscommands.ICmdObj, err error) {
				assert.Error(t, err)
				assert.EqualValues(t, "Must specify a remote if specifying a branch", err.Error())
			},
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.testName, func(t *testing.T) {
			instance := buildSyncCommands(commonDeps{})
			s.test(instance.PushCmdObj(s.opts))
		})
	}
}

func TestSyncFetch(t *testing.T) {
	type scenario struct {
		testName       string
		opts           FetchOptions
		fetchAllConfig bool
		test           func(oscommands.ICmdObj)
	}

	scenarios := []scenario{
		{
			testName:       "Fetch in foreground (all=false)",
			opts:           FetchOptions{Background: false},
			fetchAllConfig: false,
			test: func(cmdObj oscommands.ICmdObj) {
				assert.True(t, cmdObj.ShouldLog())
				assert.Equal(t, cmdObj.GetCredentialStrategy(), oscommands.PROMPT)
				assert.Equal(t, cmdObj.Args(), []string{"git", "fetch"})
			},
		},
		{
			testName:       "Fetch in foreground (all=true)",
			opts:           FetchOptions{Background: false},
			fetchAllConfig: true,
			test: func(cmdObj oscommands.ICmdObj) {
				assert.True(t, cmdObj.ShouldLog())
				assert.Equal(t, cmdObj.GetCredentialStrategy(), oscommands.PROMPT)
				assert.Equal(t, cmdObj.Args(), []string{"git", "fetch", "--all"})
			},
		},
		{
			testName:       "Fetch in background (all=false)",
			opts:           FetchOptions{Background: true},
			fetchAllConfig: false,
			test: func(cmdObj oscommands.ICmdObj) {
				assert.False(t, cmdObj.ShouldLog())
				assert.Equal(t, cmdObj.GetCredentialStrategy(), oscommands.FAIL)
				assert.Equal(t, cmdObj.Args(), []string{"git", "fetch"})
			},
		},
		{
			testName:       "Fetch in background (all=true)",
			opts:           FetchOptions{Background: true},
			fetchAllConfig: true,
			test: func(cmdObj oscommands.ICmdObj) {
				assert.False(t, cmdObj.ShouldLog())
				assert.Equal(t, cmdObj.GetCredentialStrategy(), oscommands.FAIL)
				assert.Equal(t, cmdObj.Args(), []string{"git", "fetch", "--all"})
			},
		},
	}

	for _, s := range scenarios {
		s := s
		t.Run(s.testName, func(t *testing.T) {
			instance := buildSyncCommands(commonDeps{})
			instance.UserConfig.Git.FetchAll = s.fetchAllConfig
			s.test(instance.FetchCmdObj(s.opts))
		})
	}
}
