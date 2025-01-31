package telegram

func (p *Processor) RegisterHandlers() {
	p.Handle("/start", p.handleStart)
}

func (p *Processor) handleStart() error {
	// p.tg.Send("Hello", p.tg.ChatID)
	return nil
}
