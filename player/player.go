package player

import (
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	_ "golang.org/x/oauth2"
	"os"
	_ "wrappinator.auth"
	auth "wrappinator.auth"
)

const (
	redirectURL = "http://localhost:8080/callback"

	//actions
	interrupting_playback    string = "interrupting_playback"
	pausing                  string = "pausing"
	resuming                 string = "resuming"
	seeking                  string = "seeking"
	skipping_next            string = "skipping_next"
	skipping_prev            string = "skipping_prev"
	toggle_repeat_context    string = "toggle_repeat_context"
	toggling_shuffle         string = "toggling_shuffle"
	toggling_repeat_playback string = "toggling_repeat_playback"
	transferring_playback    string = "transferring_playback"
)

var (
	envloaderr                = godotenv.Load()
	clientId           string = os.Getenv("CLIENT_ID")
	clientSecret       string = os.Getenv("CLIENT_SECRET")
	tokenStorePath     string = os.Getenv("TOKEN_STORE_PATH")
	tokenStoreFileName string = os.Getenv("TOKEN_FILENAME")
	state              string = "abc123"
	conf                      = auth.New(auth.WithRedirectURL(redirectURL), auth.WithClientID(clientId), auth.WithClientSecret(clientSecret), auth.WithScopes(auth.ScopeUserReadPrivate))
	validToken         oauth2.Token
)

type PlayerOpt func(player *Player)

func withActions(actions []string) PlayerOpt {
	return func(p *Player) {
		p.actions = actions
	}
}
func withDevice(device Device) PlayerOpt {
	return func(p *Player) {
		p.device = device
	}
}

func withRepeat_state(repeat_state string) PlayerOpt {
	return func(p *Player) {
		p.repeat_state = repeat_state
	}
}

func withShuffle_state(shuffle_state string) PlayerOpt {
	return func(p *Player) {
		p.shuffle_state = shuffle_state

	}
}
func withContext(context string) PlayerOpt {
	return func(p *Player) {
		p.context = context
	}
}
func withTimestamp(timestamp string) PlayerOpt {
	return func(p *Player) {
		p.timestamp = timestamp
	}
}
func withProgress_ms(progress_ms string) PlayerOpt {
	return func(p *Player) {
		p.progress_ms = progress_ms
	}
}
func withIs_playing(is_playing string) PlayerOpt {
	return func(p *Player) {
		p.is_playing = is_playing
	}
}
func withItem(item string) PlayerOpt {
	return func(p *Player) {
		p.item = item

	}
}
func withCurrently_playing_Type(currently_playing_type string) PlayerOpt {
	return func(p *Player) {
		p.currently_playing_type = currently_playing_type
	}
}

type Player struct {
	device                 Device
	repeat_state           string
	shuffle_state          string
	context                string
	timestamp              string
	progress_ms            string
	is_playing             string
	item                   string
	currently_playing_type string
	actions                []string
}

func New(playerOpts ...PlayerOpt) *Player {
	p := &Player{}
	for _, opt := range playerOpts {
		opt(p)
	}
	return p
}
