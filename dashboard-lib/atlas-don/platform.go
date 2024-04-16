package atlas_don

import "fmt"

type PlatformOpts struct {
	// Platform is infrastructure deployment platform: docker or k8s
	Platform     string
	LabelFilters map[string]string
	LabelFilter  string
	LegendString string
	LabelQuery   string
}

// PlatformPanelOpts generate different queries for "docker" and "k8s" deployment platforms
func PlatformPanelOpts(platform string, ocrVersion string) PlatformOpts {
	po := PlatformOpts{
		LabelFilters: map[string]string{
			"contract": `=~"${contract}"`,
		},
	}

	variableFeedId := "feed_id"
	if ocrVersion == "ocr3" {
		variableFeedId = "feed_id_name"
	}

	switch ocrVersion {
	case "ocr":
		break
	case "ocr2":
		po.LabelFilters[variableFeedId] = `=~"${` + variableFeedId + `}"`
		break
	case "ocr3":
		po.LabelFilters[variableFeedId] = `=~"${` + variableFeedId + `}"`
		break
	}
	switch platform {
	case "kubernetes":
		po.LabelFilters["namespace"] = `=~"${namespace}"`
		po.LabelFilters["job"] = `=~"${job}"`
		po.LabelFilters["pod"] = `=~"${pod}"`
		po.LabelFilter = "job"
		po.LegendString = "pod"
		break
	case "docker":
		po.LabelFilters["instance"] = `=~"${instance}"`
		po.LabelFilter = "instance"
		po.LegendString = "instance"
		break
	default:
		panic(fmt.Sprintf("failed to generate Platform dependent queries, unknown platform: %s", platform))
	}
	for key, value := range po.LabelFilters {
		po.LabelQuery += key + value + ", "
	}
	return po
}
