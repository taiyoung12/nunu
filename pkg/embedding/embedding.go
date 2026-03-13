package embedding

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
	"github.com/spf13/viper"
)

// Client wraps OpenAI embedding API.
type Client struct {
	client     *openai.Client
	model      openai.EmbeddingModel
	dimensions int
}

func NewClient(conf *viper.Viper) *Client {
	apiKey := conf.GetString("embedding.api_key")
	model := conf.GetString("embedding.model")
	dimensions := conf.GetInt("embedding.dimensions")
	if dimensions <= 0 {
		dimensions = 1536
	}

	var embModel openai.EmbeddingModel
	switch model {
	case "text-embedding-3-large":
		embModel = openai.LargeEmbedding3
	default:
		embModel = openai.SmallEmbedding3
	}

	return &Client{
		client:     openai.NewClient(apiKey),
		model:      embModel,
		dimensions: dimensions,
	}
}

// Embed generates an embedding vector for the given text.
func (c *Client) Embed(ctx context.Context, text string) ([]float32, error) {
	resp, err := c.client.CreateEmbeddings(ctx, openai.EmbeddingRequestStrings{
		Input:      []string{text},
		Model:      c.model,
		Dimensions: c.dimensions,
	})
	if err != nil {
		return nil, fmt.Errorf("embedding error: %w", err)
	}
	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("no embedding returned")
	}
	return resp.Data[0].Embedding, nil
}
