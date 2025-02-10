package service

func DefaultSelectArgs() map[string]string {
	return map[string]string{
		"per_page": "10",
		"page":     "1",
		"order":    "ASC",
		"order_by": "created_at",
	}
}

func ParseSelectQueryArgs(queryArgs map[string]string) (args map[string]string) {
	args = DefaultSelectArgs()
	for key, value := range queryArgs {
		if key != "" && value != "" {
			if _, ok := args[key]; ok {
				args[key] = value
			}
		}
	}
	return args
}
