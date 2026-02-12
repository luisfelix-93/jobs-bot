package mongodb

import (
	"context"
	"fmt"
	"time"

	"jobs-bot/internal/domain"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type aiAnalysisDoc struct {
	Score          int      `bson:"score"`
	Strengths      []string `bson:"strengths"`
	Gaps           []string `bson:"gaps"`
	Recommendation string   `bson:"recommendation"`
	Summary        string   `bson:"summary"`
	Source         string   `bson:"source"`
}

type processedJobDoc struct {
	GUID            string         `bson:"guid"`
	Source          string         `bson:"source"`
	Profile         string         `bson:"profile"`
	Title           string         `bson:"title"`
	Link            string         `bson:"link"`
	Location        string         `bson:"location"`
	Description     string         `bson:"description"`
	MatchPercentage float64        `bson:"match_percentage"`
	FoundKeywords   []string       `bson:"found_keywords"`
	MissingKeywords []string       `bson:"missing_keywords"`
	AIAnalysis      *aiAnalysisDoc `bson:"ai_analysis,omitempty"`
	Notified        bool           `bson:"notified"`
	NotifiedAt      time.Time      `bson:"notified_at"`
	CreatedAt       time.Time      `bson:"created_at"`
	TTLExpireAt     time.Time      `bson:"ttl_expire_at"`
}

type MongoJobStore struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoJobStore(uri, dbName string) (*MongoJobStore, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar ao MongoDB: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("erro ao fazer ping no MongoDB: %w", err)
	}

	collection := client.Database(dbName).Collection("processed_jobs")

	store := &MongoJobStore{
		client:     client,
		collection: collection,
	}

	if err := store.ensureIndexes(ctx); err != nil {
		return nil, fmt.Errorf("erro ao criar indexes: %w", err)
	}

	return store, nil
}

func (s *MongoJobStore) ensureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "guid", Value: 1}, {Key: "profile", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "ttl_expire_at", Value: 1}},
			Options: options.Index().SetExpireAfterSeconds(0),
		},
	}

	_, err := s.collection.Indexes().CreateMany(ctx, indexes)
	return err
}

func (s *MongoJobStore) Exists(guid, profile string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{Key: "guid", Value: guid}, {Key: "profile", Value: profile}}

	var result bson.M
	err := s.collection.FindOne(ctx, filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("erro ao verificar existência do job: %w", err)
	}

	return true, nil
}

func (s *MongoJobStore) Save(job domain.ProcessedJob) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var aiDoc *aiAnalysisDoc
	if job.AIAnalysis != nil {
		aiDoc = &aiAnalysisDoc{
			Score:          job.AIAnalysis.Score,
			Strengths:      job.AIAnalysis.Strengths,
			Gaps:           job.AIAnalysis.Gaps,
			Recommendation: job.AIAnalysis.Recommendation,
			Summary:        job.AIAnalysis.Summary,
			Source:         job.AIAnalysis.Source,
		}
	}

	doc := processedJobDoc{
		GUID:            job.GUID,
		Source:          job.Source,
		Profile:         job.Profile,
		Title:           job.Title,
		Link:            job.Link,
		Location:        job.Location,
		Description:     job.Description,
		MatchPercentage: job.KeywordAnalysis.MatchPercentage,
		FoundKeywords:   job.KeywordAnalysis.FoundKeywords,
		MissingKeywords: job.KeywordAnalysis.MissingKeywords,
		AIAnalysis:      aiDoc,
		Notified:        job.Notified,
		NotifiedAt:      job.NotifiedAt,
		CreatedAt:       job.CreatedAt,
		TTLExpireAt:     job.TTLExpireAt,
	}

	_, err := s.collection.InsertOne(ctx, doc)
	if err != nil {
		return fmt.Errorf("erro ao salvar job no MongoDB: %w", err)
	}

	return nil
}

func (s *MongoJobStore) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.client.Disconnect(ctx)
}
