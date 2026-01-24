package thetalake

var (
	mediaTypeNameToId = map[string]int{
		"video":      1,
		"audio":      2,
		"chat":       3,
		"attachment": 4,
		"email":      5,
		"image":      6,
	}

	mediaTypeIdToName = map[int]string{
		1: "video",
		2: "audio",
		3: "chat",
		4: "attachment",
		5: "email",
		6: "image",
	}
)

type MediaType struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func MediaTypesNamesToIds(types []string) []int64 {
	var ids []int64 = make([]int64, 0, len(types))
	for _, t := range types {
		id, ok := mediaTypeNameToId[t]
		if !ok {
			return nil
		}
		ids = append(ids, int64(id))
	}

	return ids
}

func MediaTypeIdsToNames(ids []int64) []string {
	var names []string = make([]string, 0, len(ids))
	for _, id := range ids {
		name, ok := mediaTypeIdToName[int(id)]
		if !ok {
			return nil
		}
		names = append(names, name)
	}

	return names
}
