package main

var videos []Video = []Video{
	{Id: "1", Title: "Video 1", Thumbnail: "thumbnail 1", Source: "source 1", ExternalId: "externalId 1"},
	{Id: "2", Title: "Video 2", Thumbnail: "thumbnail 2", Source: "source 2", ExternalId: "externalId 2"},
	{Id: "3", Title: "Video 3", Thumbnail: "thumbnail 3", Source: "source 3", ExternalId: "externalId 3"},
	{Id: "4", Title: "Video 4", Thumbnail: "thumbnail 4", Source: "source 4", ExternalId: "externalId 4"},
	{Id: "5", Title: "Video 5", Thumbnail: "thumbnail 5", Source: "source 5", ExternalId: "externalId 5"},
	{Id: "6", Title: "Video 6", Thumbnail: "thumbnail 6", Source: "source 6", ExternalId: "externalId 6"},
	{Id: "7", Title: "Video 7", Thumbnail: "thumbnail 7", Source: "source 7", ExternalId: "externalId 7"},
	{Id: "8", Title: "Video 8", Thumbnail: "thumbnail 8", Source: "source 8", ExternalId: "externalId 8"},
	{Id: "9", Title: "Video 9", Thumbnail: "thumbnail 9", Source: "source 9", ExternalId: "externalId 9"},
}

type videoRepository struct {
}

func NewVideoRepository() *videoRepository {
	return &videoRepository{}
}

func (r *videoRepository) FindAll() ([]Video, error) {
	return videos, nil
}
