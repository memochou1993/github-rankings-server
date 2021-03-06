package operator

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Unwind(field string) bson.D {
	return bson.D{
		{"$unwind", fmt.Sprintf("$%s", field)},
	}
}

func Match(field string, v interface{}) bson.D {
	return bson.D{
		{"$match", bson.D{
			{field, v},
		}},
	}
}

func Project(v interface{}) bson.D {
	return bson.D{
		{"$project", v},
	}
}

func First(field string) bson.D {
	return bson.D{
		{"$first", field},
	}
}

func In(v interface{}) bson.D {
	return bson.D{
		{"$in", v},
	}
}

func Sum(field string) bson.D {
	return bson.D{
		{"$sum", fmt.Sprintf("$%s", field)},
	}
}

func Group(v interface{}) bson.D {
	return bson.D{
		{"$group", v},
	}
}

func Sort(field string, order int64) bson.D {
	return bson.D{
		{"$sort", bson.D{
			{field, order},
		}},
	}
}

func Count() bson.D {
	return bson.D{
		{"$count", "count"},
	}
}

func Skip(skip int64) bson.D {
	return bson.D{
		{"$skip", skip},
	}
}

func Limit(limit int64) bson.D {
	return bson.D{
		{"$limit", limit},
	}
}

func Regex(pattern, options string) bson.D {
	return bson.D{
		{"$regex", primitive.Regex{Pattern: pattern, Options: options}},
	}
}
