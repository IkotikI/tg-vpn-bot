# VPN Telegram Bot

The application provide managment of VPN (3x-ui) subscriptions and serving it over Telegram Bot. Admin
managment is carried out through web admin panel (Go+HTMX).

## Installation 

### For Development
Clone the repo \
`git clone ...`\
Install development go packages \
`go install`\
`go install github.com/a-h/templ/cmd/templ@latest`\
`go install github.com/air-verse/air@latest`\
Install npm packages for Tailwindcss\
`cd ./web/admin_panel/ && npm install`

## Run

Provide enviroment variables, example of which are given in `template.env` file.

### Run admin panel (with Go air)
`make admin` \
Open provided host in web browser.

### Run Telegram Bot (with Go air)
To run bot for Telegram you need provide a valid token. ([See official guide](https://core.telegram.org/bots/tutorial))
`make tg` \
Open the Telegram Bot, which Token been provided.

## Deployment

Telegram Bot and Admin Web Panel can run on a sigle server as monolith.\
Meanwhile, there can be unlimited VPN servers, that could be deployed on different locations. One of it can be at the same server with contoll app, though.

To add new server to VPN network, you need firstly to
1. Deploy [3x-ui panel](https://github.com/MHSanaei/3x-ui) on a server. It can be do with docker-compose (see /web/3x-ui/docker/)
2. Then, get your authorization data (login, passord) from 3x-ui panel.
3. Then open admin web panel of the app, navigate to Servers -> Add new server. Create new server with given credentials.
4. [Unimplemented] For automatic configuration press "Configure"

