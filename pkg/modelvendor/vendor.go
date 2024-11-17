package modelvendor

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/llms/openai"
)

const (
	Ollama string = "ollama"
	OpenAI string = "openai"
)

type Config struct {
	VendorName string `json:"vendor_name" yaml:"vendor_name" mapstructure:"vendor_name"`
	ModelName  string `json:"model_name" yaml:"model_name" mapstructure:"model_name"`
}

type LLM interface {
	Call(ctx context.Context, prompt string, options ...llms.CallOption) (string, error)
	CreateEmbedding(ctx context.Context, texts []string) ([][]float32, error)
}

// NOTE(Hue): Returning an interface while not generally a good idea, but in this case,
// it's a good idea because we can have multiple implementations of the LLM interface.
func NewLLM(config Config) (LLM, error) {
	switch config.VendorName {
	case Ollama:
		v, err := newOllama(config)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize new Ollama llm: %w", err)
		}
		return v, nil
	case OpenAI:
		v, err := newOpenAI(config)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize new OpenAI llm: %w", err)
		}
		return v, nil
	default:
		return nil, fmt.Errorf("unknown vendor: %s", config.VendorName)
	}
}

func newOllama(config Config) (*ollama.LLM, error) {
	llm, err := ollama.New(
		ollama.WithModel(config.ModelName),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Ollama client: %w", err)
	}
	return llm, nil
}

func newOpenAI(config Config) (*openai.LLM, error) {
	llm, err := openai.New(
		openai.WithModel(config.ModelName),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OpenAI client: %w", err)
	}
	return llm, nil
}
