package mongo

import (
	"reflect"

	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mo "go.mongodb.org/mongo-driver/mongo"
	moOpts "go.mongodb.org/mongo-driver/mongo/options"
)

func deref(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

//Find finds the documents matching a model.
func (conn *connection) Find(filter interface{}, outputVal interface{}, opts ...*Options) (err error) {
	var Opts *moOpts.FindOptions
	var cur *mo.Cursor

	value := reflect.ValueOf(outputVal)
	if value.Kind() != reflect.Ptr {
		return
	}
	direct := reflect.Indirect(value)
	slice := deref(value.Type())
	if slice.Kind() != reflect.Slice {
		return
	}
	isPtr := slice.Elem().Kind() == reflect.Ptr
	base := deref(slice.Elem())

	if opts != nil {
		Opts = new(moOpts.FindOptions)
		err = copier.Copy(Opts, opts[0])
		if err != nil {
			return err
		}

		if opts[0].CollectionName == "" {
			cur, err = conn.collection.Find(conn.ctx, filter, Opts)
		} else {
			cur, err = conn.multiCollection[opts[0].CollectionName].Find(conn.ctx, filter, Opts)
		}
	} else {
		cur, err = conn.collection.Find(conn.ctx, filter)
	}

	if err != nil {
		return
	}

	for cur.Next(conn.ctx) {
		vp := reflect.New(base)

		// Create a value into which the single document can be decoded
		err := cur.Decode(vp.Interface())
		if err != nil {
			return err
		}

		if isPtr {
			direct.Set(reflect.Append(direct, vp))
		} else {
			direct.Set(reflect.Append(direct, reflect.Indirect(vp)))
		}

	}

	if err := cur.Err(); err != nil {
		cur.Close(conn.ctx)
		return err
	}

	// Close the cursor once finished
	cur.Close(conn.ctx)

	return
}

//FindOne returns up to one document that matches the model.
func (conn *connection) FindOne(filter interface{}, outputVal interface{}, opts ...*Options) error {
	var Opts *moOpts.FindOneOptions
	var res *mo.SingleResult
	if opts != nil {
		Opts = new(moOpts.FindOneOptions)
		copier.Copy(Opts, opts[0])

		if opts[0].CollectionName == "" {
			res = conn.collection.FindOne(conn.ctx, filter, Opts)
		} else {
			res = conn.multiCollection[opts[0].CollectionName].FindOne(conn.ctx, filter, Opts)
		}
	} else {
		res = conn.collection.FindOne(conn.ctx, filter)
	}
	if res.Err() != nil {
		return res.Err()
	}

	if err := res.Decode(outputVal); err != nil {
		return err
	}

	return nil
}

//FindOneAndUpdate finds a single document and updates it, returning either the original or the updated.
func (conn *connection) FindOneAndUpdate(filter, update interface{}, outputVal interface{}, opts ...*Options) error {
	var Opts *moOpts.FindOneAndUpdateOptions
	var res *mo.SingleResult
	if opts != nil {
		Opts = new(moOpts.FindOneAndUpdateOptions)
		err := copier.Copy(Opts, opts[0])
		if err != nil {
			return err
		}

		if opts[0].CollectionName == "" {
			res = conn.collection.FindOneAndUpdate(conn.ctx, filter, update, Opts)
		} else {
			res = conn.multiCollection[opts[0].CollectionName].FindOneAndUpdate(conn.ctx, filter, update, Opts)
		}
	} else {
		res = conn.collection.FindOneAndUpdate(conn.ctx, filter, update)
	}

	if res.Err() != nil {
		return res.Err()
	}

	if err := res.Decode(outputVal); err != nil {
		return err
	}

	return nil
}

//InsertOne inserts a single document into the collection
func (conn *connection) InsertOne(document interface{}, opts ...*Options) (insertedID string, err error) {
	var result *mo.InsertOneResult
	if opts != nil && opts[0].CollectionName != "" {
		result, err = conn.multiCollection[opts[0].CollectionName].InsertOne(conn.ctx, document)
	} else {
		result, err = conn.collection.InsertOne(conn.ctx, document)
	}
	if err != nil {
		return
	}

	insertedID = result.InsertedID.(primitive.ObjectID).Hex()

	return
}

//DeleteOne deletes a single document from the collection.
func (conn *connection) DeleteOne(filter interface{}, opts ...*Options) (DeletedCount int64, err error) {
	var result *mo.DeleteResult
	if opts != nil && opts[0].CollectionName != "" {
		result, err = conn.multiCollection[opts[0].CollectionName].DeleteOne(conn.ctx, filter)
	} else {
		result, err = conn.collection.DeleteOne(conn.ctx, filter)
	}
	if err != nil {
		return
	}
	DeletedCount = result.DeletedCount
	return
}

//DeleteMany deletes multiple documents from the collection.
func (conn *connection) DeleteMany(filter interface{}, opts ...*Options) (DeletedCount int64, err error) {
	var result *mo.DeleteResult
	if opts != nil && opts[0].CollectionName != "" {
		result, err = conn.multiCollection[opts[0].CollectionName].DeleteMany(conn.ctx, filter)
	} else {
		result, err = conn.collection.DeleteMany(conn.ctx, filter)
	}
	if err != nil {
		return
	}
	DeletedCount = result.DeletedCount
	return
}

//CountDocuments gets the number of documents matching the filter.
func (conn *connection) CountDocuments(filter interface{}, opts ...*Options) (total int64, err error) {
	var Opts *moOpts.CountOptions
	if opts != nil {
		Opts = new(moOpts.CountOptions)
		copier.Copy(Opts, opts[0])
		if opts[0].Skip == 0 {
			Opts.Skip = nil
		}
		if opts[0].Limit == 0 {
			Opts.Limit = nil
		}

		if opts[0].CollectionName == "" {
			total, err = conn.collection.CountDocuments(conn.ctx, filter, Opts)
		} else {
			total, err = conn.multiCollection[opts[0].CollectionName].CountDocuments(conn.ctx, filter, Opts)
		}
	} else {
		total, err = conn.collection.CountDocuments(conn.ctx, filter)
	}

	if err != nil {
		return
	}

	return
}

//Aggregate runs an aggregation framework pipeline.
func (conn *connection) Aggregate(pipeline interface{}, outputVal interface{}, opts ...*Options) (err error) {
	var cur *mo.Cursor
	value := reflect.ValueOf(outputVal)
	if value.Kind() != reflect.Ptr {
		err = ErrOutputValNotPointer
		return
	}
	direct := reflect.Indirect(value)
	slice := deref(value.Type())
	if slice.Kind() != reflect.Slice {
		err = ErrOutputValNotSlicePointer
		return
	}
	isPtr := slice.Elem().Kind() == reflect.Ptr
	base := deref(slice.Elem())

	if opts != nil && opts[0].CollectionName != "" {
		cur, err = conn.multiCollection[opts[0].CollectionName].Aggregate(conn.ctx, pipeline)
	} else {
		cur, err = conn.collection.Aggregate(conn.ctx, pipeline)
	}
	if err != nil {
		return
	}

	for cur.Next(conn.ctx) {
		vp := reflect.New(base)

		// Create a value into which the single document can be decoded
		err := cur.Decode(vp.Interface())
		if err != nil {
			return err
		}

		if isPtr {
			direct.Set(reflect.Append(direct, vp))
		} else {
			direct.Set(reflect.Append(direct, reflect.Indirect(vp)))
		}
	}

	if err := cur.Err(); err != nil {
		cur.Close(conn.ctx)
		return err
	}

	// Close the cursor once finished
	cur.Close(conn.ctx)
	return

}

//UpdateMany updates multiple documents in the collection.
func (conn *connection) UpdateMany(filter, update interface{}, opts ...*Options) (modifiedCount int64, err error) {
	var result *mo.UpdateResult
	if opts != nil && opts[0].CollectionName != "" {
		result, err = conn.multiCollection[opts[0].CollectionName].UpdateMany(conn.ctx, filter, update)
	} else {
		result, err = conn.collection.UpdateMany(conn.ctx, filter, update)
	}
	if err != nil {
		return
	}
	modifiedCount = result.ModifiedCount

	return
}

//UpdateOne updates a single document in the collection.
func (conn *connection) UpdateOne(filter, update interface{}, opts ...*Options) (upsertedID string, modifiedCount int64, err error) {
	var result *mo.UpdateResult

	var Opts *moOpts.UpdateOptions
	if opts != nil {
		Opts = new(moOpts.UpdateOptions)
		copier.Copy(Opts, opts[0])

		if opts[0].CollectionName == "" {
			result, err = conn.collection.UpdateOne(conn.ctx, filter, update, Opts)
		} else {
			result, err = conn.multiCollection[opts[0].CollectionName].UpdateOne(conn.ctx, filter, update, Opts)
		}
	} else {
		result, err = conn.collection.UpdateOne(conn.ctx, filter, update)
	}

	if err != nil {
		return
	}
	modifiedCount = result.ModifiedCount
	if result.UpsertedID != nil {
		upsertedID = result.UpsertedID.(primitive.ObjectID).Hex()
	}
	return
}
