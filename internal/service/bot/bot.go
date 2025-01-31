package bot

type Fetcher interface {
	Fetch()
}

type Processor interface {
	Process()
}

type Consumer interface {
	Start()
}
