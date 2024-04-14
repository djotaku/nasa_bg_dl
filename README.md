# NASA_BG_DL

Downloads the images from NASA's background of the day. Feed at https://www.nasa.gov/feeds/iotd-feed

It then separates the images into directories based on whether it's widescreen/landscape, square(ish), or tall/portrait. This way you can easily tell if you should use it on your computer desktop in normal orientation or if you have some of your screens rotated/want to use it on your phone. 

It will grab the 3 most recent images because sometimes NASA posts more than one image per day. 


## To use

At $HOME/.config/nasa_bg_dl/settings.json fill in the following

```json
{
    "tmp":"path to tmp directory",
    "wide":"path to landscape, or widescreen, directory",
    "square":"path to 4:3 to 3:4 direectory",
    "tall":"path to a portrait directory"
}
```

Then either run it once per day or put it in a cron job to run once per day.

logs will go to $HOME/.local/share/nasa_bg_dl/nasa_bg_dl.log

## TODO
- Make sure errors also end up in the log
- use Lumberjack to rotate logs