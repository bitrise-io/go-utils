package git

import (
	"strconv"
	"strings"

	"github.com/bitrise-io/go-utils/command"
)

// Init creates an empty Git reporitory or reinitializes an existing one.
func (g *Git) Init() *command.Model {
	return g.Command("init")
}

// Clone a repository into a new directory
func (g *Git) Clone(repo string) *command.Model {
	return g.Command("clone", repo)
}

// RemoteList shows a list of existing remote urls with remote names
func (g *Git) RemoteList() *command.Model {
	return g.Command("remote", "-v")
}

// RemoteAdd adds a remote named <name> for the repository at <url>.
func (g *Git) RemoteAdd(name, url string) *command.Model {
	return g.Command("remote", "add", name, url)
}

// Fetch downloads objects and refs from another repository.
func (g *Git) Fetch(depth int, opts ...string) *command.Model {
	args := []string{"fetch"}
	if depth != 0 {
		args = append(args, "--depth="+strconv.Itoa(depth))
	}
	args = append(args, opts...)
	return g.Command(args...)
}

// FetchPR downloads objects and refs for the pull request branch.
func (g *Git) FetchPR(branch string) *command.Model {
	return g.Command("fetch", "origin", branch+":"+strings.TrimSuffix(branch, "/merge"))
}

// Checkout switchs branches or restore working tree files.
// Arg can be a commit hash, a branch or a tag.
func (g *Git) Checkout(arg string) *command.Model {
	return g.Command("checkout", arg)
}

// Merge joins two or more development histories together.
// Arg can be a commit hash, branch or tag.
func (g *Git) Merge(arg string) *command.Model {
	return g.Command("merge", arg)
}

// Reset the current branch head to commit and possibly updates the index.
// The mode must be one of the following: --soft, --mixed, --hard, --merge, --keep.
func (g *Git) Reset(mode, commit string) *command.Model {
	return g.Command("reset", mode, commit)
}

// Clean removes untracked files from the working tree.
func (g *Git) Clean(options ...string) *command.Model {
	args := []string{"clean"}
	args = append(args, options...)
	return g.Command(args...)
}

// SubmoduleUpdate updates the registered submodules.
func (g *Git) SubmoduleUpdate() *command.Model {
	return g.Command("submodule", "update", "--init", "--recursive")
}

// SubmoduleForeach evaluates an arbitrary git command in each checked out
// submodule.
func (g *Git) SubmoduleForeach(command *command.Model) *command.Model {
	args := []string{"submodule", "foreach"}
	args = append(args, command.GetCmd().Args...)
	return g.Command(args...)
}

// Pull incorporates changes from a remote repository into the current branch.
func (g *Git) Pull() *command.Model {
	return g.Command("pull")
}

// Add file contents to the index
func (g *Git) Add(pathspec string) *command.Model {
	return g.Command("add", pathspec)
}

// Branch lists branches
func (g *Git) Branch() *command.Model {
	return g.Command("branch")
}

// NewBranch creates a new branch as if git-branch were called and then check it out.
func (g *Git) NewBranch(branch string) *command.Model {
	return g.Command("checkout", "-b", branch)
}

// Apply reads the supplied diff output (patch) and applies it to files.
func (g *Git) Apply(patch string) *command.Model {
	return g.Command("apply", "--index", patch)
}

// Log shows the commit logs. The format parameter controls what is shown and how.
func (g *Git) Log(format string) *command.Model {
	return g.Command("log", "-1", "--format="+format)
}

// RevList lists commit objects in reverse chronological order.
func (g *Git) RevList(commit string, opts ...string) *command.Model {
	args := []string{"rev-list", commit}
	args = append(args, opts...)
	return g.Command(args...)
}

// Push updates remote refs along with associated objects.
func (g *Git) Push(branch string) *command.Model {
	return g.Command("push", "-u", "origin", branch)
}

// Commit Stores the current contents of the index in a new commit along with a log message from the user describing the changes.
func (g *Git) Commit(message string) *command.Model {
	return g.Command("commit", "-m", message)
}

// RevParse picks out and massage parameters
func (g *Git) RevParse(arg string) *command.Model {
	return g.Command("rev-parse", arg)
}

// Status shows the working tree status.
func (g *Git) Status(opts ...string) *command.Model {
	args := []string{"status"}
	args = append(args, opts...)
	return g.Command(args...)
}

// CloneTagOrBranch is recursively clones a tag or branch
func (g *Git) CloneTagOrBranch(uri, destination, tagOrBranch string) *command.Model {
	return g.Command("git", "clone", "--recursive", "--branch", tagOrBranch, uri, destination)
}
