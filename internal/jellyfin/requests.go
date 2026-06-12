package jellyfin

import (
	"context"

	"github.com/sj14/jellyfin-go/api"
)

// itemFields are the fields requested for items in list views.
// Genres is needed for genre tags in descriptions; MediaStreams for video info.
var itemFields = []api.ItemFields{api.ITEMFIELDS_MEDIA_STREAMS, api.ITEMFIELDS_GENRES}

func (c *Client) GetResume() ([]Item, error) {
	res, _, err := c.api.ItemsAPI.GetResumeItems(context.Background()).
		UserId(c.UserID).
		Fields(itemFields).
		Execute()
	if err != nil {
		return nil, err
	}
	if res.Items == nil {
		return []Item{}, nil
	}
	return res.Items, nil
}

func (c *Client) GetNextUp() ([]Item, error) {
	res, _, err := c.api.TvShowsAPI.GetNextUp(context.Background()).
		Fields(itemFields).
		EnableTotalRecordCount(false).
		DisableFirstEpisode(false).
		EnableResumable(false).
		EnableRewatching(false).
		Execute()
	if err != nil {
		return nil, err
	}
	if res.Items == nil {
		return []Item{}, nil
	}
	return res.Items, nil
}

func (c *Client) GetRecentlyAdded() ([]Item, error) {
	res, _, err := c.api.ItemsAPI.GetItems(context.Background()).
		Recursive(true).
		IncludeItemTypes([]api.BaseItemKind{api.BASEITEMKIND_MOVIE, api.BASEITEMKIND_SERIES}).
		Fields(itemFields).
		Limit(100).
		SortBy([]api.ItemSortBy{api.ITEMSORTBY_DATE_CREATED}).
		SortOrder([]api.SortOrder{api.SORTORDER_DESCENDING}).
		Execute()
	if err != nil {
		return nil, err
	}
	if res.Items == nil {
		return []Item{}, nil
	}
	return res.Items, nil
}

func (c *Client) GetLibraries() ([]Item, error) {
	res, _, err := c.api.UserViewsAPI.GetUserViews(context.Background()).
		UserId(c.UserID).
		Execute()
	if err != nil {
		return nil, err
	}
	if res.Items == nil {
		return []Item{}, nil
	}
	return res.Items, nil
}

func (c *Client) GetSeasons(seriesId string) ([]Item, error) {
	res, _, err := c.api.TvShowsAPI.GetSeasons(context.Background(), seriesId).
		UserId(c.UserID).
		Fields(itemFields).
		Execute()
	if err != nil {
		return nil, err
	}
	if res.Items == nil {
		return []Item{}, nil
	}
	return res.Items, nil
}

func (c *Client) GetEpisodes(seriesId, seasonId string) ([]Item, error) {
	req := c.api.TvShowsAPI.GetEpisodes(context.Background(), seriesId).
		Fields(itemFields)
	if seasonId != "" {
		req = req.SeasonId(seasonId)
	}
	res, _, err := req.Execute()
	if err != nil {
		return nil, err
	}
	if res.Items == nil {
		return []Item{}, nil
	}
	return res.Items, nil
}

func (c *Client) GetItemsByParent(parentId string) ([]Item, error) {
	res, _, err := c.api.ItemsAPI.GetItems(context.Background()).
		ParentId(parentId).
		Fields(itemFields).
		Execute()
	if err != nil {
		return nil, err
	}
	if res.Items == nil {
		return []Item{}, nil
	}
	return res.Items, nil
}

func (c *Client) Search(query string) ([]Item, error) {
	res, _, err := c.api.ItemsAPI.GetItems(context.Background()).
		SearchTerm(query).
		Recursive(true).
		IncludeItemTypes([]api.BaseItemKind{api.BASEITEMKIND_MOVIE, api.BASEITEMKIND_SERIES}).
		Fields(itemFields).
		Limit(100).
		Execute()
	if err != nil {
		return nil, err
	}
	if res.Items == nil {
		return []Item{}, nil
	}
	return res.Items, nil
}

func (c *Client) ReportPlaybackStart(item Item, ticks int64) error {
	_, err := c.api.PlaystateAPI.ReportPlaybackStart(context.Background()).PlaybackStartInfo(api.PlaybackStartInfo{
		ItemId:        item.Id,
		PositionTicks: *api.NewNullableInt64(&ticks),
	}).Execute()
	return err
}

func (c *Client) ReportPlaybackStopped(item Item, ticks int64) error {
	_, err := c.api.PlaystateAPI.ReportPlaybackStopped(context.Background()).PlaybackStopInfo(api.PlaybackStopInfo{
		ItemId:        item.Id,
		PositionTicks: *api.NewNullableInt64(&ticks),
	}).Execute()
	return err
}

func (c *Client) ReportPlaybackProgress(item Item, ticks int64) error {
	_, err := c.api.PlaystateAPI.ReportPlaybackProgress(context.Background()).PlaybackProgressInfo(api.PlaybackProgressInfo{
		ItemId:        item.Id,
		PositionTicks: *api.NewNullableInt64(&ticks),
	}).Execute()
	return err
}

// GetMediaSegments returns a map of start ticks to end ticks of media segments
//
//   - item: the item to get media segments for
//   - types: array of media segment types to include. If empty, returns nil.
func (c *Client) GetMediaSegments(item Item, types []string) (map[int64]int64, error) {
	if len(types) == 0 {
		return nil, nil
	}
	// cast []string to []api.MediaSegmentType
	mediaSegmentTypes := make([]api.MediaSegmentType, len(types))
	for i, t := range types {
		mediaSegmentTypes[i] = api.MediaSegmentType(t)
	}
	res, _, err := c.api.MediaSegmentsAPI.GetItemSegments(context.Background(), item.GetId()).IncludeSegmentTypes(mediaSegmentTypes).Execute()
	if err != nil {
		return nil, err
	}
	segments := make(map[int64]int64, len(res.Items))
	for _, segment := range res.Items {
		segments[segment.GetStartTicks()] = segment.GetEndTicks()
	}
	return segments, nil
}

// GetAllMediaSegments returns a map of start ticks to end ticks of all media segments for the given item.
func (c *Client) GetAllMediaSegments(item Item) (map[int64]int64, error) {
	res, _, err := c.api.MediaSegmentsAPI.GetItemSegments(context.Background(), item.GetId()).Execute()
	if err != nil {
		return nil, err
	}
	segments := make(map[int64]int64, len(res.Items))
	for _, segment := range res.Items {
		segments[segment.GetStartTicks()] = segment.GetEndTicks()
	}
	return segments, nil
}

func (c *Client) MarkAsWatched(item Item) error {
	_, _, err := c.api.PlaystateAPI.MarkPlayedItem(context.Background(), item.GetId()).Execute()
	return err
}

func (c *Client) MarkAsUnwatched(item Item) error {
	_, _, err := c.api.PlaystateAPI.MarkUnplayedItem(context.Background(), item.GetId()).Execute()
	return err
}
