package ctrl

import (
	"fmt"
)

// WriteChannel save channel
func WriteChannel(uuid string, channel *Channel) *Channel {
	globalCtrl.hubLock.Lock()
	data, ok := globalCtrl.channelHub[uuid]
	if !ok {
		data = new(Channel)
		data.Params = map[string]string{}
		globalCtrl.channelHub[uuid] = data
	}
	globalCtrl.hubLock.Unlock()

	data.lock.Lock()
	data.CtrlUuid = channel.CtrlUuid
	data.subs = channel.subs
	if channel.NodeUuid != "" {
		data.NodeUuid = channel.NodeUuid
	}
	if channel.Uuid != "" {
		data.Uuid = channel.Uuid
	}
	if channel.PeerUuid != "" {
		data.PeerUuid = channel.PeerUuid
	}
	// call direction inbound |outbound
	if channel.Direction != "" {
		data.Direction = channel.Direction
	}
	// START RINGING ANSWERED ACTIVE DESTROY READY ...
	if channel.State != "" {
		data.State = channel.State
	}
	if channel.CidName != "" {
		data.CidName = channel.CidName
	}
	if channel.CidNumber != "" {
		data.CidNumber = channel.CidNumber
	}
	if channel.DestNumber != "" {
		data.DestNumber = channel.DestNumber
	}

	if channel.CreateEpoch > 0 {
		data.CreateEpoch = channel.CreateEpoch
	}

	if channel.RingEpoch > 0 {
		data.RingEpoch = channel.RingEpoch
	}

	if channel.AnswerEpoch > 0 {
		data.AnswerEpoch = channel.AnswerEpoch
	}

	if channel.HangupEpoch > 0 {
		data.HangupEpoch = channel.HangupEpoch
	}

	if channel.Answered {
		data.Answered = channel.Answered

	}
	// list of uuids
	if len(channel.Peers) > 0 {
		data.Peers = channel.Peers
	}
	if channel.Params != nil {
		for k, v := range channel.Params {
			data.Params[k] = v
		}
	}
	data.lock.Unlock()
	return data
}

func FindChannel(condition string, argument string) []*Channel {
	matched := make([]*Channel, 0)
	globalCtrl.hubLock.RLock()
	for _, channel := range globalCtrl.channelHub {
		if channel.GetVariable0(condition) == argument {
			matched = append(matched, channel)
		}
	}
	globalCtrl.hubLock.RUnlock()
	return matched
}

// ReadChannel get channel
func ReadChannel(uuid string) (*Channel, error) {
	globalCtrl.hubLock.RLock()
	data, ok := globalCtrl.channelHub[uuid]
	if ok {
		globalCtrl.hubLock.RUnlock()
		return data, nil
	}
	globalCtrl.hubLock.RUnlock()
	return nil, fmt.Errorf("not found")
}

// DelChannel get channel
func DelChannel(uuid string) error {
	globalCtrl.hubLock.Lock()
	channel, ok := globalCtrl.channelHub[uuid]
	if ok {
		for _, sub := range channel.subs {
			sub.Unsubscribe()
		}
	}
	delete(globalCtrl.channelHub, uuid)
	globalCtrl.hubLock.Unlock()
	return nil
}

// GetChannelState 获取 channel 状态
func GetChannelState(uuid string) string {
	if uuid != "" {
		if channel, err := ReadChannel(uuid); err == nil {
			return channel.GetState()
		}
	}
	return ""
}
