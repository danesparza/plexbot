# dlshow [![CircleCI](https://circleci.com/gh/danesparza/dlshow.svg?style=shield)](https://circleci.com/gh/danesparza/dlshow)
Go package to help parse downloaded TV show filenames

## Example

```golang
filename := "Once.Upon.a.Time.S03E01.720p.HDTV.X264-DIMENSION.mkv"
showInfo, err := dlshow.GetEpisodeInfo(filename)

// showInfo.ShowName == "Once Upon a Time"
// showInfo.SeasonNumber == 3
// showInfo.EpisodeNumber == 1
```
