package thetalake

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
)

func TestCreateLabel(t *testing.T) {
	labelRequest := Label{
		BackgroundColor: "#FFC906",
		Hidden:          false,
		LongName:        "This is a test label",
		ShortName:       "Test Label",
	}

	testHandler := func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		var receivedLabel Label
		err := json.Unmarshal(body, &receivedLabel)
		assert.NoError(t, err)

		assert.Equal(t, labelRequest.BackgroundColor, receivedLabel.BackgroundColor)
		assert.Equal(t, labelRequest.Hidden, receivedLabel.Hidden)
		assert.Equal(t, labelRequest.LongName, receivedLabel.LongName)
		assert.Equal(t, labelRequest.ShortName, receivedLabel.ShortName)

		w.WriteHeader(http.StatusCreated)
		w.Write(readTestData("create_label_response.json"))
	}

	client := newTestClient(t, http.MethodPost, "/labels", testHandler)

	createdLabel, err := client.CreateLabel(context.Background(), labelRequest)
	assert.NoError(t, err)

	assert.Equal(t, int64(5), createdLabel.Id)
	assert.Equal(t, "Test Label", createdLabel.ShortName)
	assert.Equal(t, "This is a test label", createdLabel.LongName)
	assert.Equal(t, "#FFC906", createdLabel.BackgroundColor)
	assert.Equal(t, false, createdLabel.Hidden)
}

func TestGetLabelById(t *testing.T) {
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(readTestData("get_label_by_id_response.json"))
	}

	client := newTestClient(t, http.MethodGet, "/labels/5", testHandler)

	retrievedLabel, err := client.GetLabelById(context.Background(), 5)
	assert.NoError(t, err)

	assert.Equal(t, int64(5), retrievedLabel.Id)
	assert.Equal(t, "Test Label", retrievedLabel.ShortName)
	assert.Equal(t, "This is a test label", retrievedLabel.LongName)
	assert.Equal(t, "#FFC906", retrievedLabel.BackgroundColor)
	assert.Equal(t, false, retrievedLabel.Hidden)
	assert.Equal(t, int64(7), retrievedLabel.TaggedDatumsCount)
}

func TestUpdateLabel(t *testing.T) {
	labelUpdateRequest := Label{
		Id:              5,
		BackgroundColor: "#FF0000",
		Hidden:          true,
		LongName:        "Updated test label",
		ShortName:       "Updated Label",
	}

	testHandler := func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		var receivedLabel Label
		err := json.Unmarshal(body, &receivedLabel)
		assert.NoError(t, err)

		assert.Equal(t, labelUpdateRequest.BackgroundColor, receivedLabel.BackgroundColor)
		assert.Equal(t, labelUpdateRequest.Hidden, receivedLabel.Hidden)
		assert.Equal(t, labelUpdateRequest.LongName, receivedLabel.LongName)
		assert.Equal(t, labelUpdateRequest.ShortName, receivedLabel.ShortName)

		w.WriteHeader(http.StatusOK)
		w.Write(readTestData("update_label_response.json"))
	}

	client := newTestClient(t, http.MethodPut, "/labels/5", testHandler)

	updatedLabel, err := client.UpdateLabel(context.Background(), labelUpdateRequest)
	assert.NoError(t, err)

	assert.Equal(t, int64(5), updatedLabel.Id)
	assert.Equal(t, "Updated Label", updatedLabel.ShortName)
	assert.Equal(t, "Updated test label", updatedLabel.LongName)
	assert.Equal(t, "#FF0000", updatedLabel.BackgroundColor)
	assert.Equal(t, true, updatedLabel.Hidden)
}

func TestDeleteLabel(t *testing.T) {
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(readTestData("get_label_by_id_response.json"))
	}

	client := newTestClient(t, http.MethodDelete, "/labels/5", testHandler)

	err := client.DeleteLabel(context.Background(), 5)
	assert.NoError(t, err)
}
