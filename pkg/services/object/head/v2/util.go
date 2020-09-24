package headsvc

import (
	"github.com/nspcc-dev/neofs-api-go/pkg/object"
	objectV2 "github.com/nspcc-dev/neofs-api-go/v2/object"
	headsvc "github.com/nspcc-dev/neofs-node/pkg/services/object/head"
)

func toPrm(req *objectV2.HeadRequest) *headsvc.Prm {
	body := req.GetBody()

	return new(headsvc.Prm).
		WithAddress(
			object.NewAddressFromV2(body.GetAddress()),
		).
		Short(body.GetMainOnly()).
		OnlyLocal(req.GetMetaHeader().GetTTL() == 1) // FIXME: use constant
}

func fromResponse(r *headsvc.Response, short bool) *objectV2.HeadResponse {
	fn := fullPartFromResponse
	if short {
		fn = shortPartFromResponse
	}

	body := new(objectV2.HeadResponseBody)
	body.SetHeaderPart(fn(r))

	resp := new(objectV2.HeadResponse)
	resp.SetBody(body)

	return resp
}

func fullPartFromResponse(r *headsvc.Response) objectV2.GetHeaderPart {
	obj := r.Header().ToV2()

	hs := new(objectV2.HeaderWithSignature)
	hs.SetHeader(obj.GetHeader())
	hs.SetSignature(obj.GetSignature())

	p := new(objectV2.GetHeaderPartFull)
	p.SetHeaderWithSignature(hs)

	return p
}

func shortPartFromResponse(r *headsvc.Response) objectV2.GetHeaderPart {
	hdr := r.Header().ToV2().GetHeader()

	sh := new(objectV2.ShortHeader)
	sh.SetOwnerID(hdr.GetOwnerID())
	sh.SetCreationEpoch(hdr.GetCreationEpoch())
	sh.SetPayloadLength(hdr.GetPayloadLength())
	sh.SetVersion(hdr.GetVersion())
	sh.SetObjectType(hdr.GetObjectType())

	p := new(objectV2.GetHeaderPartShort)
	p.SetShortHeader(sh)

	return p
}