# goscreener



## Description

A simple service for creating screenshots of site pages, based on chromedp

## Usage

Use docker compose
```
docker-compose up
```
* Open url <u>localhost:30222</u> 

Or you can use go run
```
go run cmd/server/main.go
```
- Open url <u>localhost:8080</u> 
<br>
<br>

If the page opens on your host machine, then everything is fine


## Example
Make post requests to the **/screenshot** route and fill the body with json
```
POST /screenshot HTTP/1.1
Host: localhost:8080
Content-Type: application/json
Content-Length: 59

{
  "url": "https://en.wikipedia.org/wiki/FromSoftware"
}
```

## Example
For many screens, use **/screenshots-many**
```
POST /screenshots-many HTTP/1.1
Host: localhost:8080
Content-Type: application/json
Content-Length: 1495

{
  "urls": [
    "https://en.wikipedia.org/wiki/FromSoftware",
    "https://ru.wikipedia.org/wiki/Bloodborne",
    "https://ru.wikipedia.org/wiki/Elden_Ring"
  ]
}
```

## Params
| param                           | description                                                                                      | type        | default value | example                                                                                           |
|---------------------------------|--------------------------------------------------------------------------------------------------|-------------|---------------|---------------------------------------------------------------------------------------------------|
| url                             | The url of the page what you want to screen                                                      | string      |               |                                                                                                   |
| urls                            | Multiple getting screenshots, returns screenshot links                                           | string      |               |                                                                                                   |
| cache                           | Use with "urls" param, for temporary saving screens (screens will be available within an 1 hour) | bool        | false         |                                                                                                   |
| load_page_timeout_seconds       | Timeout after navigating to a page, but before other actions (use this if content loads slowly)  | int         | 0             |                                                                                                   |
| final_load_page_timeout_seconds | Final timeout after all actions have been completed                                              | int         | 0             |                                                                                                   |
| height                          | Viewport height                                                                                  | int         | 1920          |                                                                                                   |
| width                           | Viewport width                                                                                   | int         | 1080          |                                                                                                   |
| quality                         | Image quality                                                                                    | int         | 100           |                                                                                                   |
| full_screen                     | Takes a screenshot of the entire page                                                            | bool        | false         |                                                                                                   |
| target_selector                 | Takes a screenshot started from the top and ending at the selected target                        | string      |               | "target_selector":".endless__item"                                                                |
| fake_nav                        | Imitation an additional browser-style panel located at the top                                   | bool        | false         |                                                                                                   |
| with_scroll                     | Makes smooth scrolling actions of the page, first down to the bottom and then to the top         | bool        | false         |                                                                                                   |
| remove_nodes                    | Removes nodes from the DOM                                                                       | string json |               | "remove_nodes": [{"selector":"div[data-type=\"banner\"], div.banner","parent":false,"many":true}] |
| fixed_nodes                     | Fixed positions of elements in the DOM                                                           | string json |               | "fixed_nodes" : [{"selector":"#headerSticked"]}]                                                  |
