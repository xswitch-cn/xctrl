package tboy

type TBoySimple struct {
	*TBoy
}

// NewTBoy .
func NewSimple(node_uuid string, domain string, fn ...OptionFn) *TBoySimple {
	boy := &TBoySimple{
		TBoy: &TBoy{
			Options: &Options{},
		},
	}
	boy.Init()
	boy.SetUUID(node_uuid)
	boy.SetDomain(domain)
	for _, o := range fn {
		o(boy.Options)
	}

	log.Infof("new boy created uuid=%s domain=%s", boy.NodeUUID, boy.Domain)
	return boy
}

func (boy *TBoySimple) AddChannel(key string, channel *FakeChannel) {
	Channels[key] = channel
}
