package loop

const KeepAliveTickDuration = keepAliveTickDuration

func (p *RelayerService) Kill() {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if p.client != nil {
		p.client.Kill()
	}
}

func (p *RelayerService) Reset() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.client != nil {
		p.client.Kill()
	}
	p.client = nil
	p.clientProtocol = nil
	p.plug = nil
	p.relayer = nil
}
