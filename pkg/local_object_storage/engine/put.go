package engine

import (
	"errors"

	"github.com/nspcc-dev/neofs-node/pkg/core/object"
	"github.com/nspcc-dev/neofs-node/pkg/local_object_storage/shard"
	"go.uber.org/zap"
)

// PutPrm groups the parameters of Put operation.
type PutPrm struct {
	obj *object.Object
}

// PutRes groups resulting values of Put operation.
type PutRes struct{}

var errPutShard = errors.New("could not put object to any shard")

// WithObject is a Put option to set object to save.
//
// Option is required.
func (p *PutPrm) WithObject(obj *object.Object) *PutPrm {
	if p != nil {
		p.obj = obj
	}

	return p
}

// Put saves the object to local storage.
//
// Returns any error encountered that
// did not allow to completely save the object.
func (e *StorageEngine) Put(prm *PutPrm) (*PutRes, error) {
	if e.metrics != nil {
		defer elapsed(e.metrics.AddPutDuration)()
	}

	_, err := e.exists(prm.obj.Address()) // todo: make this check parallel
	if err != nil {
		return nil, err
	}

	existPrm := new(shard.ExistsPrm)
	existPrm.WithAddress(prm.obj.Address())

	finished := false

	e.iterateOverSortedShards(prm.obj.Address(), func(ind int, s *shard.Shard) (stop bool) {
		exists, err := s.Exists(existPrm)
		if err != nil {
			return false // this is not ErrAlreadyRemoved error so we can go to the next shard
		}

		if exists.Exists() {
			if ind != 0 {
				toMoveItPrm := new(shard.ToMoveItPrm)
				toMoveItPrm.WithAddress(prm.obj.Address())

				_, err = s.ToMoveIt(toMoveItPrm)
				if err != nil {
					e.log.Warn("could not mark object for shard relocation",
						zap.Stringer("shard", s.ID()),
						zap.String("error", err.Error()),
					)
				}
			}

			finished = true

			return true
		}

		putPrm := new(shard.PutPrm)
		putPrm.WithObject(prm.obj)

		_, err = s.Put(putPrm)
		if err != nil {
			e.log.Warn("could not put object in shard",
				zap.Stringer("shard", s.ID()),
				zap.String("error", err.Error()),
			)

			return false
		}

		finished = true

		return true
	})

	if !finished {
		err = errPutShard
	}

	return nil, err
}

// Put writes provided object to local storage.
func Put(storage *StorageEngine, obj *object.Object) error {
	_, err := storage.Put(new(PutPrm).
		WithObject(obj),
	)

	return err
}
