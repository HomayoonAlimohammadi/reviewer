package main

import (
	"context"
	"fmt"
	"log"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores/milvus"
)

type Caller interface {
	Call(ctx context.Context, prompt string, options ...llms.CallOption) (string, error)
}

type Reviewer struct {
	caller      Caller
	milvusStore milvus.Store
}

func main() {
	log.Println("starting...")

	// llmClient, err := openai.New(openai.WithToken(os.Getenv("OPENAI_API_KEY")))
	llmClient, err := ollama.New(ollama.WithModel("llama3.2:3b"))
	if err != nil {
		log.Fatalln("failed to create openai client:", err)
	}

	log.Println("llm client created")

	embedder, err := embeddings.NewEmbedder(llmClient)
	if err != nil {
		log.Fatalln("failed to create embedder:", err)
	}

	log.Println("embedder created")

	milvusStore, err := initMilvusStore(context.Background(), embedder)
	if err != nil {
		log.Fatalln("failed to initialize milvus store:", err)
	}

	log.Println("milvus store initialized")

	rv := &Reviewer{
		caller:      llmClient,
		milvusStore: milvusStore,
	}

	guidelines, err := retrieveGuideLines(context.Background())
	if err != nil {
		log.Fatalln("failed to retrieve guidelines:", err)
	}

	log.Println("guidelines retrieved")

	if err := rv.StoreBestPractices(context.Background(), milvusStore, embedder, guidelines); err != nil {
		log.Fatalln("failed to store best practices:", err)
	}

	log.Println("best practices stored")

	sampleCode := `
	GLOBAL_VAR = 5
	def bad_name_that_does_not_mean_anything():
		a_bad_variable_name = 10
		GLOBAL_VAR += 1
		print(a_bad_variable_name)
	`
	docs, err := rv.RetrieveTopGuidelines(context.Background(), sampleCode)
	if err != nil {
		log.Fatalln("failed to retrieve best practices:", err)
	}

	log.Println("top guidelines retrieved")

	feedback, err := rv.GenerateFeedback(context.Background(), docs, sampleCode)
	if err != nil {
		log.Fatalln("failed to generate feedback:", err)
	}

	log.Println("feedback generated")

	fmt.Println(feedback)

	log.Println("done!")
}

func (r *Reviewer) StoreBestPractices(ctx context.Context, store milvus.Store, embedder embeddings.Embedder, data []string) error {
	for _, text := range data {
		_, err := store.AddDocuments(ctx,
			[]schema.Document{{PageContent: text}},
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Reviewer) RetrieveTopGuidelines(ctx context.Context, codeSnippet string) ([]schema.Document, error) {
	docs, err := r.milvusStore.SimilaritySearch(ctx, codeSnippet, 5)
	if err != nil {
		return nil, fmt.Errorf("failed to do similarity search: %w", err)
	}

	return docs, nil
}

func (r *Reviewer) GenerateFeedback(ctx context.Context, retrievedData []schema.Document, codeSnippet string) (string, error) {
	guidelines := ""
	for _, doc := range retrievedData {
		guidelines += doc.PageContent + "\n"
	}

	prompt := fmt.Sprintf(`
	Review this code:
	%s
	Based on these guidelines:
	%s
	`,
		codeSnippet,
		guidelines,
	)

	feedback, err := r.caller.Call(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("failed to call LLM: %w", err)
	}

	return feedback, nil
}

func initMilvusStore(ctx context.Context, embedder embeddings.Embedder) (milvus.Store, error) {
	index, err := entity.NewIndexIvfFlat(entity.L2, 16384)
	if err != nil {
		return milvus.Store{}, fmt.Errorf("failed to create new index: %w", err)
	}

	client, err := milvus.New(ctx,
		client.Config{
			Address: "localhost:19530",
		},
		milvus.WithEmbedder(embedder),
		milvus.WithIndex(index),
	)
	if err != nil {
		return milvus.Store{}, fmt.Errorf("failed to initialize new milvus: %w", err)
	}

	return client, nil
}

func retrieveGuideLines(_ context.Context) ([]string, error) {
	return []string{
		"Use meaningful variable names",
		"Use comments to explain complex code",
		"Keep functions short and focused",
		"Use consistent indentation",
		"Avoid global variables",
	}, nil
}
