package mpris

import (
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

	getPropertyMethod = "org.freedesktop.DBus.Properties.Get"
	setPropertyMethod = "org.freedesktop.DBus.Properties.Set"
)

func setProperty(obj *dbus.Object, iface string, prop string, val interface{}) error {
	call := obj.Call(setPropertyMethod, 0, prop, val)
	return call.Err
}

func List(conn *dbus.Conn) ([]string, error) {
	var names []string
	err := conn.BusObject().Call("org.freedesktop.DBus.ListNames", 0).Store(&names)
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
	return prop.String(), nil
}

type player struct {
	obj *dbus.Object
}

func (i *player) Next() {
	i.obj.Call(playerInterface+".Next", 0)
}

func (i *player) Previous() {
	i.obj.Call(playerInterface+".Previous", 0)
}

func (i *player) Pause() {
	i.obj.Call(playerInterface+".Pause", 0)
}

func (i *player) PlayPause() {
	i.obj.Call(playerInterface+".PlayPause", 0)
}

func (i *player) Stop() {
	i.obj.Call(playerInterface+".Stop", 0)
}

func (i *player) Play() {
	i.obj.Call(playerInterface+".Play", 0)
}

func (i *player) Seek(offset int64) {
	i.obj.Call(playerInterface+".Seek", 0, offset)
}

func (i *player) SetPosition(trackId *dbus.ObjectPath, position int64) {
	i.obj.Call(playerInterface+".SetPosition", 0, trackId, position)
}

func (i *player) OpenUri(uri string) {
	i.obj.Call(playerInterface+".OpenUri", 0, uri)
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
	return PlaybackStatus(variant.String()), nil
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
	return LoopStatus(prop.String()), nil
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
	return setProperty(i.obj, playerInterface, "Volume", volume)
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
