package model

import (
	"context"
	"fmt"
	"github.com/memochou1993/github-rankings/app/handler/request"
	"github.com/memochou1993/github-rankings/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

type Rank struct {
	Name       string    `json:"name" bson:"name"`
	ImageUrl   string    `json:"imageUrl" bson:"image_url"`
	Rank       int       `json:"rank" bson:"rank"`
	Last       int       `json:"last" bson:"last"`
	TotalCount int       `json:"totalCount" bson:"total_count"`
	Tags       []string  `json:"tags" bson:"tags"`
	CreatedAt  time.Time `json:"createdAt" bson:"created_at"`
}

type RankRecord struct {
	ID         string `bson:"_id"`
	ImageUrl   string `bson:"image_url"`
	TotalCount int    `bson:"total_count"`
}

type RankModel struct {
	*Model
}

func (r *RankModel) CreateIndexes() {
	database.CreateIndexes(r.Name(), []string{
		"name",
		"tags",
	})
}

func (r *RankModel) List(req *request.Request, timestamps []time.Time) []Rank {
	ctx := context.Background()
	cond := mongo.Pipeline{bson.D{{"created_at", bson.D{{"$in", timestamps}}}}}
	if !req.IsNameEmpty() {
		cond = append(cond, bson.D{{"name", req.Name}})
	}
	var tags []string
	if !req.IsTypeEmpty() {
		tags = append(tags, fmt.Sprintf("type:%s", req.Type))
	}
	if !req.IsFieldEmpty() {
		tags = append(tags, fmt.Sprintf("field:%s", req.Field))
	}
	cond = append(cond, bson.D{{"tags", tags}})
	pipeline := mongo.Pipeline{
		bson.D{
			{"$match", bson.D{
				{"$and", cond},
			}},
		},
		bson.D{
			{"$skip", (req.Page - 1) * req.Limit},
		},
		bson.D{
			{"$limit", req.Limit},
		},
	}
	cursor := database.Aggregate(ctx, NewRankModel().Name(), pipeline)
	ranks := make([]Rank, req.Limit)
	if err := cursor.All(ctx, &ranks); err != nil {
		log.Fatal(err.Error())
	}
	return ranks
}

func (r *RankModel) Store(model Interface, p Pipeline, createdAt time.Time) {
	ctx := context.Background()
	cursor := database.Aggregate(ctx, model.Name(), *p.Pipeline)
	defer database.CloseCursor(ctx, cursor)

	last := p.Count(model)

	var models []mongo.WriteModel
	for i := 0; cursor.Next(ctx); i++ {
		rec := RankRecord{}
		if err := cursor.Decode(&rec); err != nil {
			log.Fatal(err.Error())
		}

		doc := Rank{
			Name:       rec.ID,
			ImageUrl:   rec.ImageUrl,
			Rank:       i + 1,
			Last:       last,
			TotalCount: rec.TotalCount,
			Tags:       p.Tags,
			CreatedAt:  createdAt,
		}
		models = append(models, mongo.NewInsertOneModel().SetDocument(doc))
		if cursor.RemainingBatchLength() == 0 {
			database.BulkWrite(r.Name(), models)
			models = models[:0]
		}
	}
}

func (r *RankModel) Delete(createdAt time.Time, tags ...string) {
	filter := bson.D{
		{"$and", []bson.D{{
			{"tags", bson.D{
				{"$in", tags},
			}},
			{"created_at", bson.D{
				{"$lt", createdAt},
			}},
		}}},
	}
	database.DeleteMany(r.Name(), filter)
}

type Pipeline struct {
	Pipeline *mongo.Pipeline
	Tags     []string
}

func (p *Pipeline) Count(model Interface) int {
	ctx := context.Background()
	rec := struct {
		Count int `bson:"count"`
	}{}
	pipeline := append(*p.Pipeline, bson.D{{"$match", bson.D{{"total_count", bson.D{{"$gt", 0}}}}}})
	pipeline = append(pipeline, bson.D{{"$count", "count"}})
	cursor := database.Aggregate(ctx, model.Name(), pipeline)
	defer database.CloseCursor(ctx, cursor)
	for cursor.Next(ctx) {
		if err := cursor.Decode(&rec); err != nil {
			log.Fatal(err.Error())
		}
	}
	return rec.Count
}

func NewRankModel() *RankModel {
	return &RankModel{
		&Model{
			name: "ranks",
		},
	}
}
