package mocks

import (
	"fmt"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/kubeflow/model-registry/pkg/openapi"
)

func GenerateMockRegisteredModelList() openapi.RegisteredModelList {
	var models []openapi.RegisteredModel
	for i := 0; i < 2; i++ {
		model := GenerateRegisteredModel()
		models = append(models, model)
	}

	return openapi.RegisteredModelList{
		NextPageToken: gofakeit.UUID(),
		PageSize:      int32(gofakeit.Number(1, 20)),
		Size:          int32(len(models)),
		Items:         models,
	}
}

func GenerateRegisteredModel() openapi.RegisteredModel {
	model := openapi.RegisteredModel{
		CustomProperties: &map[string]openapi.MetadataValue{
			"example_key": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue:  gofakeit.Sentence(3),
					MetadataType: "string",
				},
			},
		},
		Description:              stringToPointer(gofakeit.Sentence(5)),
		ExternalId:               stringToPointer(gofakeit.UUID()),
		Name:                     stringToPointer(gofakeit.Name()),
		Id:                       stringToPointer(gofakeit.UUID()),
		CreateTimeSinceEpoch:     stringToPointer(fmt.Sprintf("%d", gofakeit.Date().UnixMilli())),
		LastUpdateTimeSinceEpoch: stringToPointer(fmt.Sprintf("%d", gofakeit.Date().UnixMilli())),
		Owner:                    stringToPointer(gofakeit.Name()),
		State:                    stateToPointer(openapi.RegisteredModelState(gofakeit.RandomString([]string{string(openapi.REGISTEREDMODELSTATE_LIVE), string(openapi.REGISTEREDMODELSTATE_ARCHIVED)}))),
	}
	return model
}

func stateToPointer(s openapi.RegisteredModelState) *openapi.RegisteredModelState {
	return &s
}

func stringToPointer(s string) *string {
	return &s
}
