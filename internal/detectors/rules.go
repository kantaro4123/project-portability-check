package detectors

// Rule describes one stable finding ID exposed through JSON, SARIF,
// configuration suppressions, and --list-rules.
type Rule struct {
	ID      string
	Summary string
}

// Rules returns the public rule catalog in stable category order.
func Rules() []Rule {
	return []Rule{
		{ID: "paths.absolute", Summary: "Machine-specific absolute user paths"},
		{ID: "fs.windows-reserved", Summary: "Windows reserved device names"},
		{ID: "fs.windows-trailing", Summary: "Windows-unsafe trailing spaces or dots"},
		{ID: "fs.windows-forbidden-char", Summary: "Windows-forbidden path characters"},
		{ID: "fs.windows-long-path", Summary: "Risky long checkout paths on Windows"},
		{ID: "fs.case-collision", Summary: "Case-insensitive filesystem collisions"},
		{ID: "fs.symlink", Summary: "Symlink portability hazards"},
		{ID: "fs.script-not-executable", Summary: "Shebang scripts missing executable permission"},
		{ID: "text.mixed-line-endings", Summary: "Mixed newline conventions"},
		{ID: "text.shell-crlf", Summary: "Executable shebang lines stored with CRLF"},
		{ID: "text.non-utf8", Summary: "Text files that are not valid UTF-8"},
		{ID: "text.utf8-bom", Summary: "UTF-8 byte order marks"},
		{ID: "shell.grep-p", Summary: "GNU-only grep -P usage"},
		{ID: "shell.sed-i", Summary: "GNU/BSD sed -i incompatibility"},
		{ID: "shell.readlink-f", Summary: "GNU-only readlink -f usage"},
		{ID: "shell.date-d", Summary: "GNU/BSD date flag incompatibility"},
		{ID: "shell.xargs-r", Summary: "GNU-only xargs -r usage"},
		{ID: "shell.sh-bashism", Summary: "Bash-only syntax under a POSIX sh shebang"},
		{ID: "imports.case-mismatch", Summary: "Relative JavaScript/TypeScript imports use incorrect path case"},
		{ID: "package.script-unix", Summary: "Unix-specific package.json scripts"},
		{ID: "runtime.node-unpinned", Summary: "Node.js runtime is not pinned"},
		{ID: "runtime.python-unpinned", Summary: "Python runtime is not pinned"},
		{ID: "runtime.go-unpinned", Summary: "go.mod lacks a Go language version"},
		{ID: "deps.node-lockfile", Summary: "JavaScript dependencies are not locked"},
		{ID: "deps.cargo-lockfile", Summary: "Cargo dependencies are not locked"},
		{ID: "env.no-example", Summary: "Environment-dependent code lacks an example env file"},
		{ID: "binary.native", Summary: "Checked-in native binary artifacts"},
		{ID: "docker.fixed-platform", Summary: "Dockerfile pins a single CPU platform"},
		{ID: "git.no-gitattributes", Summary: "Platform-sensitive scripts lack line-ending policy"},
		{ID: "ci.platform-coverage", Summary: "CI does not cover all major desktop OS families"},
	}
}
