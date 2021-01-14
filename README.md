# D-Sender

Telegram bot that can extract media resources from 2ch.hk

## Capabilities

Options for users:
* List all available origins: `/list`
* List your subscriptions: `/subs`
* Subscribe to origin: `/subscribe [origin_number]`
* Unsubscribe from origin: `/rm [subscribtion_number]`
* Create origin visible to you: `/create [board] [recource_type] [tags]`

Options for admins:
* List all available origins with description: `/clist`
* Create origin visible to everyone `/create_default [board] [recource_type] [tags] [display_name]`
* Remove origin visible to everyone `/rm_default [origin_number]`

---
## Creating origins

[board] - board name, without "/"

[resource_type] must be a string like `"( .img | .gif | webm )"`. For example, valid string is `.img.gif`
* `.img` will match image formats
* `.gif` will match gif format
* `.webm` will match video formats

[tags] must be in format `[ ! ]"tag1"{ & | | }...`
* `&` - means conjunction
* `|` - means disjunction
* `!` - means negation

Matching is made using disjunctive normal form, i.e. first it calculates negation, than conjunction, disjunction is the last\
For example, string `"cats"|"dogs"&!"big"` will match threads, description of which contains "cats" or "dogs", but not "big"

[display_name] is a string that will be visible to everyone when they call `/list`

Example: `/create_default wp .img "wallpaper"&"desktop" Wallpapers`

---
## Configuring

Disk is used to save videos in .webm format to satisfy ffmpeg requirements (I use converting because telegram does not support this format). After sending, files will be deleted automatically.

In `configs/config.yml`:
* db - database configuration
* dapi - 2ch api, you can change it to use other mirrors or custom api
* tg.admin_id - list of admins telegram id
* disk:
  * path - relative or absolute path of directory, where files will be saved
  * size - max allowed space in bytes. Files, that extends this parameter, will be discarded
* polling - period of time in minutes, after which new threads will fetched

Environment variables:
* DB_PASSWORD - database password
* DB_PORT - database port
* BOT_TOKEN - your telegram bot token

## Running

Run: `docker-compose up`



