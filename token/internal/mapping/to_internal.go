package mapping

import (
	desc "github.com/de1phin/iam/genproto/services/token/api"

	"github.com/de1phin/iam/token/internal/model"
)

func MapTokenToInternal(pb *desc.Token) model.Token {
	if pb == nil {
		return model.Token{}
	}
	return model.Token{
		Token: pb.Token,
	}
}
