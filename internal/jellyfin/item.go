package jellyfin

import (
	"fmt"
	"strings"

	"github.com/sj14/jellyfin-go/api"
)

// Helpers to try to contain BaseItemDto implementation from leaking out of the package without casting it to something else
// NOTE: just cast it to something else?

// Item is a type alias just because it looks nicer
type Item = api.BaseItemDto

func getItemRuntime(ticks int64) string {
	minutes := ticks / 600_000_000
	if minutes < 60 {
		return fmt.Sprintf("%dm", minutes)
	}
	hours := minutes / 60
	minutes -= hours * 60
	return fmt.Sprintf("%dh%dm", hours, minutes)
}

func GetResumePosition(item Item) (ticks int64) {
	if data, ok := item.GetUserDataOk(); ok {
		ticks = data.GetPlaybackPositionTicks()
	}
	return
}

func GetStreamingURL(host string, item Item) string {
	url := fmt.Sprintf("%s/videos/%s/stream?static=true", host, *item.Id)
	return fmt.Sprintf("edl://%%%d%%%s", len(url), url)
}

func GetMediaTitle(item Item) string {
	title := item.GetPath()
	switch item.GetType() {
	case api.BASEITEMKIND_MOVIE:
		title = fmt.Sprintf("%s (%d)", item.GetName(), item.GetProductionYear())
	case api.BASEITEMKIND_EPISODE:
		title = fmt.Sprintf("%s - S%d:E%d - %s (%d)", item.GetSeriesName(), item.GetParentIndexNumber(), item.GetIndexNumber(), item.GetName(), item.GetProductionYear())
	}
	return title
}

func GetItemTitle(item Item) string {
	str := &strings.Builder{}
	switch item.GetType() {
	case api.BASEITEMKIND_MOVIE:
		name := item.GetName()
		year := item.GetProductionYear()
		if year > 0 {
			fmt.Fprintf(str, "%s (%d)", name, year)
		} else {
			fmt.Fprintf(str, "%s", name)
		}
		if data, ok := item.GetUserDataOk(); ok && data.GetPlayedPercentage() > 0 {
			fmt.Fprintf(str, " [%.f%%]", data.GetPlayedPercentage())
		}
	case api.BASEITEMKIND_EPISODE:
		fmt.Fprintf(str, "%s S%.2dE%.2d (%d)", item.GetSeriesName(), item.GetParentIndexNumber(), item.GetIndexNumber(), item.GetProductionYear())
		if data, ok := item.GetUserDataOk(); ok && data.GetPlayedPercentage() > 0 {
			fmt.Fprintf(str, " [%.f%%]", data.GetPlayedPercentage())
		}
	case api.BASEITEMKIND_SERIES:
		name := item.GetName()
		year := item.GetProductionYear()
		if year > 0 {
			fmt.Fprintf(str, "%s (%d)", name, year)
		} else {
			fmt.Fprintf(str, "%s", name)
		}
		if data, ok := item.GetUserDataOk(); ok {
			fmt.Fprintf(str, " [%d]", data.GetUnplayedItemCount())
		}
	case api.BASEITEMKIND_VIDEO:
		fmt.Fprintf(str, "%s (%d)", item.GetName(), item.GetProductionYear())
	case api.BASEITEMKIND_USER_VIEW, api.BASEITEMKIND_COLLECTION_FOLDER, api.BASEITEMKIND_AGGREGATE_FOLDER, api.BASEITEMKIND_FOLDER:
		fmt.Fprintf(str, "%s", item.GetName())
		if data, ok := item.GetUserDataOk(); ok {
			unplayed := data.GetUnplayedItemCount()
			if unplayed > 0 {
				fmt.Fprintf(str, " [%d]", unplayed)
			}
		}
	}
	return str.String()
}

func GetItemDescription(item Item) string {
	str := &strings.Builder{}
	switch item.GetType() {
	case api.BASEITEMKIND_MOVIE:
		rating := item.GetCommunityRating()
		runtime := getItemRuntime(item.GetRunTimeTicks())
		if rating > 0 {
			fmt.Fprintf(str, "Movie  | Rating: %.1f | Runtime: %s", rating, runtime)
		} else {
			fmt.Fprintf(str, "Movie  | Runtime: %s", runtime)
		}
	case api.BASEITEMKIND_SERIES:
		fmt.Fprintf(str, "Series")
	case api.BASEITEMKIND_EPISODE:
		fmt.Fprintf(str, "%s", item.GetName())
	case api.BASEITEMKIND_VIDEO:
		rating := item.GetCommunityRating()
		runtime := getItemRuntime(item.GetRunTimeTicks())
		if rating > 0 {
			fmt.Fprintf(str, "Video  | Rating: %.1f | Runtime: %s", rating, runtime)
		} else {
			fmt.Fprintf(str, "Video  | Runtime: %s", runtime)
		}
	case api.BASEITEMKIND_USER_VIEW, api.BASEITEMKIND_COLLECTION_FOLDER, api.BASEITEMKIND_AGGREGATE_FOLDER, api.BASEITEMKIND_FOLDER:
		fmt.Fprintf(str, "Library")
	}
	return str.String()
}

func IsMovie(item Item) bool {
	return item.GetType() == api.BASEITEMKIND_MOVIE
}

func IsSeries(item Item) bool {
	return item.GetType() == api.BASEITEMKIND_SERIES
}

func IsEpisode(item Item) bool {
	return item.GetType() == api.BASEITEMKIND_EPISODE
}

func IsVideo(item Item) bool {
	return item.GetType() == api.BASEITEMKIND_VIDEO
}

func IsFolder(item Item) bool {
	kind := item.GetType()
	return kind == api.BASEITEMKIND_FOLDER ||
		kind == api.BASEITEMKIND_USER_VIEW ||
		kind == api.BASEITEMKIND_COLLECTION_FOLDER ||
		kind == api.BASEITEMKIND_AGGREGATE_FOLDER
}

func Watched(item Item) bool {
	if data, ok := item.GetUserDataOk(); ok {
		return data.GetPlayed()
	}
	return false
}

// ExternalSubtitleStream represents an external subtitle stream
type ExternalSubtitleStream struct {
	Language string
	Title    string
	Path     string
}

// GetExternalSubtitleStreams returns all external subtitle streams for an item
func GetExternalSubtitleStreams(item Item) []ExternalSubtitleStream {
	var subtitles []ExternalSubtitleStream
	streams := item.GetMediaStreams()
	for _, stream := range streams {
		if stream.GetType() == "Subtitle" && stream.GetIsExternal() {
			index := stream.GetIndex()
			subtitle := ExternalSubtitleStream{}
			if lang, ok := stream.GetLanguageOk(); ok && lang != nil {
				subtitle.Language = *lang
			}
			if title, ok := stream.GetDisplayTitleOk(); ok && title != nil {
				subtitle.Title = *title
			} else {
				subtitle.Title = fmt.Sprintf("External %d", index)
			}
			subtitle.Path = fmt.Sprintf("/Videos/%s/%s/Subtitles/%d/0/Stream.srt", item.GetId(), item.GetId(), index)
			subtitles = append(subtitles, subtitle)
		}
	}
	return subtitles
}
