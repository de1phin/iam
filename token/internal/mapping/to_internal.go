package mapping

import (
	desc "github.com/de1phin/iam/genproto/services/token/api"

	"github.com/de1phin/iam/token/internal/model"
)

func mapTokenToInternal(pb *desc.Token) model.Token {
	return model.Token{
		Token: pb.Token,
	}
}
