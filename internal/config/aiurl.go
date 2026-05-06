package config

type AiConfig struct {
	URL         string
	DownloadURL string
	MMOURL      string
	MMOToken    string
}

func (c *Config) InitAiURLConfig() *AiConfig {
	return &AiConfig{
		URL:         c.Ai.URL,
		DownloadURL: c.Ai.DownloadURL,
		MMOURL:      c.Ai.MMOURL,
		MMOToken:    c.Ai.MMOToken,
	}
}
