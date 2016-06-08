package incoming

import (
	"strconv"

	"golang.org/x/net/context"

	"github.com/pkg/errors"
	"google.golang.org/cloud/datastore"
)

func NewGDatastoreStorage(projectID string) *GDatastoreStorage {
	return &GDatastoreStorage{
		ProjectID: projectID,
	}
}

func (s *GDatastoreStorage) Save(ctx context.Context, events ...*ReceivedEvent) error {
	cl, err := datastore.NewClient(ctx, s.ProjectID)
	if err != nil {
		return errors.Wrap(err, "failed to create datastore client")
	}

	// classify entries into basetime / event name
	for _, e := range events {
		id := e.ReceivedOn().Unix()
		if mod := id % 300; mod > 0 {
			id = id - mod
		}

		parent := datastore.NewKey(ctx, "ReceivedEvents", strconv.FormatInt(id, 10), 0, nil)
		key := datastore.NewIncompleteKey(ctx, e.Name(), parent)
		_, err := cl.Put(ctx, key, e)
		if err != nil {
			return errors.Wrap(err, "failed to Put event to datastore")
		}
	}
	return nil
}