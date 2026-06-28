package detector

type ScanResult struct {
	Words        map[string]string
	FoundBadWord bool
	Badwords     []string
	Tokenized    []string
}
