package mpris

import (
	"context"
	"strings"

	"github.com/godbus/dbus/v5"
)

const (
	dbusObjectPath          = "/org/mpris/MediaPlayer2"
	propertiesChangedSignal = "org.freedesktop.DBus.Properties.PropertiesChanged"

	baseInterface      = "org.mpris.MediaPlayer2"
	playerInterface    = "org.mpris.MediaPlayer2.Player"
	trackListInterface = "org.mpris.MediaPlayer2.TrackList"
	playlistsInterface = "org.mpris.MediaPlayer2.Playlists"
)

func List(ctx context.Context, conn *dbus.Conn) ([]string, error) {
	var names []string
	err := conn.BusObject().CallWithContext(ctx, "org.freedesktop.DBus.ListNames", 0).Store(&names)
	if err != nil {
		return nil, err
	}

	var mprisNames []string
	prefix := baseInterface + "."
	for _, name := range names {
		if strings.HasPrefix(name, prefix) {
			mprisNames = append(mprisNames, strings.TrimPrefix(name, prefix))
		}
	}
	return mprisNames, nil
}

type Player struct {
	*base
	*player
}

type base struct {
	obj *dbus.Object
}

func (i *base) Raise() error {
	call := i.obj.Call(baseInterface+".Raise", 0)
	return call.Err
}

func (i *base) Quit() error {
	call := i.obj.Call(baseInterface+".Quit", 0)
	return call.Err
}

func (i *base) GetIdentity() (string, error) {
	prop, err := i.obj.GetProperty(baseInterface + ".Identity")
	if err != nil {
		return "", err
	}
	return prop.Value().(string), nil
}

type player struct {
	obj *dbus.Object
}

func (i *player) Next(ctx context.Context) error {
	return i.obj.CallWithContext(ctx, playerInterface+".Next", 0).Err
}

func (i *player) Previous(ctx context.Context) error {
	return i.obj.CallWithContext(ctx, playerInterface+".Previous", 0).Err
}

func (i *player) Pause(ctx context.Context) error {
	return i.obj.CallWithContext(ctx, playerInterface+".Pause", 0).Err
}

func (i *player) PlayPause(ctx context.Context) error {
	return i.obj.CallWithContext(ctx, playerInterface+".PlayPause", 0).Err
}

func (i *player) Stop(ctx context.Context) error {
	return i.obj.CallWithContext(ctx, playerInterface+".Stop", 0).Err
}

func (i *player) Play(ctx context.Context) error {
	return i.obj.CallWithContext(ctx, playerInterface+".Play", 0).Err
}

func (i *player) Seek(ctx context.Context, offset int64) error {
	return i.obj.CallWithContext(ctx, playerInterface+".Seek", 0, offset).Err
}

func (i *player) SetPosition(trackId *dbus.ObjectPath, position int64) error {
	return i.obj.Call(playerInterface+".SetPosition", 0, trackId, position).Err
}

func (i *player) OpenUri(ctx context.Context, uri string) error {
	return i.obj.CallWithContext(ctx, playerInterface+".OpenUri", 0, uri).Err
}

func (i *player) VolumeUp(ctx context.Context) error {
	return i.obj.CallWithContext(ctx, playerInterface+".VolumeUp", 0).Err
}

func (i *player) VolumeDown(ctx context.Context) error {
	return i.obj.CallWithContext(ctx, playerInterface+".VolumeDown", 0).Err
}

type PlaybackStatus string

const (
	PlaybackPlaying PlaybackStatus = "Playing"
	PlaybackPaused                 = "Paused"
	PlaybackStopped                = "Stopped"
)

func (i *player) GetPlaybackStatus() (PlaybackStatus, error) {
	variant, err := i.obj.GetProperty(playerInterface + ".PlaybackStatus")
	if err != nil {
		return "", err
	}
	return PlaybackStatus(variant.Value().(string)), nil
}

type LoopStatus string

const (
	LoopNone     LoopStatus = "None"
	LoopTrack               = "Track"
	LoopPlaylist            = "Playlist"
)

func (i *player) GetLoopStatus() (LoopStatus, error) {
	prop, err := i.obj.GetProperty(playerInterface + ".LoopStatus")
	if err != nil {
		return "", nil
	}
	return LoopStatus(prop.Value().(string)), nil
}

func (i *player) GetRate() (float64, error) {
	prop, err := i.obj.GetProperty(playerInterface + ".Rate")
	if err != nil {
		return 0, err
	}
	return prop.Value().(float64), nil
}

func (i *player) GetShuffle() (bool, error) {
	prop, err := i.obj.GetProperty(playerInterface + ".Shuffle")
	if err != nil {
		return false, err
	}
	return prop.Value().(bool), nil
}

func (i *player) GetMetadata() (map[string]dbus.Variant, error) {
	prop, err := i.obj.GetProperty(playerInterface + ".Metadata")
	if err != nil {
		return nil, err
	}
	return prop.Value().(map[string]dbus.Variant), nil
}

func (i *player) GetVolume() (float64, error) {
	prop, err := i.obj.GetProperty(playerInterface + ".Volume")
	if err != nil {
		return 0, err
	}
	return prop.Value().(float64), err
}
func (i *player) SetVolume(volume float64) error {
	return i.obj.SetProperty(playerInterface+".Volume", dbus.MakeVariant(volume))
}

func (i *player) GetPosition() (int64, error) {
	prop, err := i.obj.GetProperty(playerInterface + ".Position")
	if err != nil {
		return 0, err
	}
	return prop.Value().(int64), nil
}

func New(conn *dbus.Conn, name string) *Player {
	obj := conn.Object(baseInterface+"."+name, dbusObjectPath).(*dbus.Object)

	return &Player{
		&base{obj},
		&player{obj},
	}
}
