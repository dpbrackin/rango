package embedding

import (
	"context"
	"fmt"
	"io"
	"rango/api/internal/core"

	"github.com/openai/openai-go"
)

type OpenAIEmbedder struct {
	client *openai.Client
}

func NewOpenAIEmbedder() *OpenAIEmbedder {
	client := openai.NewClient()

	return &OpenAIEmbedder{
		client: client,
	}
}

// Embed implements core.Embedder.
func (o *OpenAIEmbedder) Embed(ctx context.Context, r io.Reader) (core.Embeddings, error) {
	contents, err := io.ReadAll(r)

	if err != nil {
		return nil, fmt.Errorf("Cannot read contents: %w", err)
	}

	contentStr := string(contents)
	inputs := openai.EmbeddingNewParamsInputArrayOfStrings([]string{contentStr})

	res, err := o.client.Embeddings.New(ctx, openai.EmbeddingNewParams{
		Input:          openai.F(openai.EmbeddingNewParamsInputUnion(inputs)),
		Model:          openai.F(openai.EmbeddingModelTextEmbedding3Small),
		EncodingFormat: openai.F(openai.EmbeddingNewParamsEncodingFormatFloat),
	})

	if err != nil {
		return nil, fmt.Errorf("Failed to get embeddings: %w", err)
	}

	embeddings := make([][]float64, len(res.Data))

	for i, e := range res.Data {
		embeddings[i] = e.Embedding
	}

	return embeddings, nil
}
