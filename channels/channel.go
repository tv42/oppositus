package channels

//go:generate go run github.com/campoy/jsonenums -type=Channel
//go:generate go run golang.org/x/tools/cmd/stringer -type=Channel

// Channel represents a CoreOS release channel.
type Channel int

// stringer doesn't understand case, so we need to make the consts be
// exactly the strings we want in JSON, and then redefine the exported
// consts in terms of the unexported ones.
// https://github.com/campoy/jsonenums/issues/13

const (
	_ Channel = iota
	stable
	beta
	alpha
)

// CoreOS release channels.
const (
	Stable = stable
	Beta   = beta
	Alpha  = alpha
)

// All returns all release channels. Callers should not mutate the
// returned data.
func All() []Channel {
	return []Channel{Stable, Beta, Alpha}
}
