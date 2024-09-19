package api

import (
	"github.com/kubeflow/model-registry/pkg/openapi"
	"github.com/kubeflow/model-registry/ui/bff/internal/mocks"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGetModelVersionHandler(t *testing.T) {
	data := mocks.GetModelVersionMocks()[0]
	expected := ModelVersionEnvelope{Data: &data}

	actual, rs, err := setupApiTest[ModelVersionEnvelope](http.MethodGet, "/api/v1/model_registry/model-registry/model_versions/1", nil)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rs.StatusCode)
	assert.Equal(t, expected.Data.Name, actual.Data.Name)
}

func TestCreateModelVersionHandler(t *testing.T) {
	data := mocks.GetModelVersionMocks()[0]
	expected := ModelVersionEnvelope{Data: &data}

	body := ModelVersionEnvelope{Data: openapi.NewModelVersion("Model One", "1")}

	actual, rs, err := setupApiTest[ModelVersionEnvelope](http.MethodPost, "/api/v1/model_registry/model-registry/model_versions", body)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusCreated, rs.StatusCode)
	assert.Equal(t, expected.Data.Name, actual.Data.Name)
	assert.Equal(t, rs.Header.Get("Location"), "/api/v1/model_registry/model-registry/model_versions/1")
}

func TestUpdateModelVersionHandler(t *testing.T) {
	data := mocks.GetModelVersionMocks()[0]
	expected := ModelVersionEnvelope{Data: &data}

	reqData := openapi.ModelVersionUpdate{
		Description: openapi.PtrString("New description"),
	}
	body := ModelVersionUpdateEnvelope{Data: &reqData}

	actual, rs, err := setupApiTest[ModelVersionEnvelope](http.MethodPatch, "/api/v1/model_registry/model-registry/model_versions/1", body)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rs.StatusCode)
	assert.Equal(t, expected.Data.Name, actual.Data.Name)
}

func TestGetAllModelArtifactsByModelVersionHandler(t *testing.T) {
	data := mocks.GetModelArtifactListMock()
	expected := ModelArtifactListEnvelope{Data: &data}

	actual, rs, err := setupApiTest[ModelArtifactListEnvelope](http.MethodGet, "/api/v1/model_registry/model-registry/model_versions/1/artifacts", nil)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rs.StatusCode)
	assert.Equal(t, expected.Data.Size, actual.Data.Size)
	assert.Equal(t, expected.Data.PageSize, actual.Data.PageSize)
	assert.Equal(t, expected.Data.NextPageToken, actual.Data.NextPageToken)
	assert.Equal(t, len(expected.Data.Items), len(actual.Data.Items))
}

func TestCreateModelArtifactByModelVersionHandler(t *testing.T) {
	data := mocks.GetModelArtifactMocks()[0]
	expected := ModelArtifactEnvelope{Data: &data}

	artifact := openapi.ModelArtifact{
		Name:         openapi.PtrString("Artifact One"),
		ArtifactType: "ARTIFACT_TYPE_ONE",
	}
	body := ModelArtifactEnvelope{Data: &artifact}

	actual, rs, err := setupApiTest[ModelArtifactEnvelope](http.MethodPost, "/api/v1/model_registry/model-registry/model_versions/1/artifacts", body)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusCreated, rs.StatusCode)
	assert.Equal(t, expected.Data.GetArtifactType(), actual.Data.GetArtifactType())
	assert.Equal(t, rs.Header.Get("Location"), "/api/v1/model_registry/model-registry/model_artifacts/1")
}
