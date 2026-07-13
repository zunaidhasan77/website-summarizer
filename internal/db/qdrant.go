package db

import (
	"context"
	"log"

	"github.com/qdrant/go-client/qdrant"
)

var QdrantClient *qdrant.Client

func InitQdrant() {
	var err error
	QdrantClient, err = qdrant.NewClient(&qdrant.Config{
		Host: "localhost",
		Port: 6334,
	})
	if err != nil {
		log.Fatalf("Could not initialize Qdrant: %v", err)
	}
}

func EnsureCollection(name string) {
	ctx := context.Background()

	err := QdrantClient.CreateCollection(ctx, &qdrant.CreateCollection{
		CollectionName: name,
		VectorsConfig: &qdrant.VectorsConfig{
			Config: &qdrant.VectorsConfig_Params{
				Params: &qdrant.VectorParams{
					Size:     768,
					Distance: qdrant.Distance_Cosine,
				},
			},
		},
	})

	if err != nil {
		log.Println("Note: Collection creation result (could be 'exists'):", err)
	}
}

func UpsertChunk(collection string, id uint64, vector []float32, text string) error {
	points := []*qdrant.PointStruct{
		{
			Id: qdrant.NewIDNum(id),
			// Reverted to the correct syntax that actually works
			Vectors: qdrant.NewVectors(vector...),
			Payload: qdrant.NewValueMap(map[string]interface{}{"text": text}),
		},
	}

	_, err := QdrantClient.Upsert(context.Background(), &qdrant.UpsertPoints{
		CollectionName: collection,
		Points:         points,
	})
	return err
}

func SearchChunks(collection string, vector []float32, limit uint64) ([]string, error) {
	ctx := context.Background()

	// 1. Build the modern QueryPoints request
	limitPtr := limit
	searchReq := &qdrant.QueryPoints{
		CollectionName: collection,
		Query:          qdrant.NewQuery(vector...),
		Limit:          &limitPtr,
		WithPayload:    qdrant.NewWithPayload(true),
	}

	// 2. Execute the search using the new Query method
	resp, err := QdrantClient.Query(ctx, searchReq)
	if err != nil {
		return nil, err
	}

	// 3. Extract the text strings from the response payload
	var results []string
	for _, point := range resp {
		if payload, ok := point.Payload["text"]; ok {
			results = append(results, payload.GetStringValue())
		}
	}

	return results, nil
}
