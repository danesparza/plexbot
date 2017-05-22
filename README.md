# plexbot [![CircleCI](https://circleci.com/gh/danesparza/plexbot.svg?style=svg)](https://circleci.com/gh/danesparza/plexbot)
Simple app to help organize tv shows and movies into the Plex naming format

# Quick start
Generate a config file:
`plexbot defaults > plexbot.yaml`

After updating your config, run the plexbot on a file:
`plexbot --config c:\plexbot\plexbot.yaml move "%F"`
where %F is the content path
