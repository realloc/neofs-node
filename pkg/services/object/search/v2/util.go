package searchsvc

import (
	"github.com/nspcc-dev/neofs-api-go/pkg/container"
	"github.com/nspcc-dev/neofs-api-go/v2/object"
	"github.com/nspcc-dev/neofs-api-go/v2/refs"
	searchsvc "github.com/nspcc-dev/neofs-node/pkg/services/object/search"
	"github.com/nspcc-dev/neofs-node/pkg/services/object/search/query"
	queryV1 "github.com/nspcc-dev/neofs-node/pkg/services/object/search/query/v1"
	"github.com/nspcc-dev/neofs-node/pkg/services/object/util"
	"github.com/pkg/errors"
)

func toPrm(body *object.SearchRequestBody, req *object.SearchRequest) (*searchsvc.Prm, error) {
	var q query.Query

	switch v := body.GetVersion(); v {
	default:
		return nil, errors.Errorf("unsupported query version #%d", v)
	case 1:
		fs := body.GetFilters()
		fsV1 := make([]*queryV1.Filter, 0, len(fs))

		for _, f := range fs {
			switch mt := f.GetMatchType(); mt {
			default:
				return nil, errors.Errorf("unsupported match type %d in query version #%d", mt, v)
			case object.MatchStringEqual:
				fsV1 = append(fsV1, queryV1.NewFilterEqual(f.GetName(), f.GetValue()))
			}
		}

		q = queryV1.New(fsV1...)
	}

	return new(searchsvc.Prm).
		WithContainerID(
			container.NewIDFromV2(body.GetContainerID()),
		).
		WithSearchQuery(q).
		WithCommonPrm(util.CommonPrmFromV2(req)), nil
}

func fromResponse(r *searchsvc.Response) *object.SearchResponse {
	ids := r.IDList()
	idsV2 := make([]*refs.ObjectID, 0, len(ids))

	for i := range ids {
		idsV2 = append(idsV2, ids[i].ToV2())
	}

	body := new(object.SearchResponseBody)
	body.SetIDList(idsV2)

	resp := new(object.SearchResponse)
	resp.SetBody(body)

	return resp
}
