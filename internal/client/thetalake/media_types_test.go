package thetalake

import (
	"testing"

	"github.com/go-openapi/testify/v2/assert"
)

func TestMediaTypesNamesToIds(t *testing.T) {

	actual := MediaTypesNamesToIds([]string{"audio", "video"})
	expected := []int64{2, 1}
	assert.Equal(t, expected, actual)

	actual = MediaTypesNamesToIds([]string{"chat", "attachment", "email"})
	expected = []int64{3, 4, 5}
	assert.Equal(t, expected, actual)

	actual = MediaTypesNamesToIds([]string{"image"})
	expected = []int64{6}
	assert.Equal(t, expected, actual)

	actual = MediaTypesNamesToIds([]string{"invalid_media_type"})
	assert.Nil(t, actual)
}

func TestMediaTypeIdsToNames(t *testing.T) {
	actual := MediaTypeIdsToNames([]int64{1, 2})
	expected := []string{"video", "audio"}
	assert.Equal(t, expected, actual)

	actual = MediaTypeIdsToNames([]int64{3, 4, 5})
	expected = []string{"chat", "attachment", "email"}
	assert.Equal(t, expected, actual)

	actual = MediaTypeIdsToNames([]int64{6})
	expected = []string{"image"}
	assert.Equal(t, expected, actual)

	actual = MediaTypeIdsToNames([]int64{999})
	assert.Nil(t, actual)
}
