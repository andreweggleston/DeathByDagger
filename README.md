# DeathByDagger
CSH annual Death by Dagger backend and slack bot. Written in Go.

## Requirements
* Go >= 1.10
* PostgreSQL

## Installation
This project requires the use of environment variables. 

| Environment Variable | Description |
|----------------------|-------------|
|    `SERVER_ADDR`     |Address to serve on|
|    `PUBLIC_ADDR`     |Publicly accessible address for the server, requires schema|
|    `ALLOWED_ORIGINS`     ||
|    `DATABASE_ADDR`     |Database Address|
|    `DATABASE_NAME`     |Database Name|
|    `DATABASE_USERNAME`     |Database Username|
|    `DATABASE_PASSWORD`     |Database password|
|    `COOKIE_STORE_SECRET`     |base64 encoded key to use for encrypting cookies|
|    `SECURE_COOKIE`     |Enable 'secure' flag on cookies|
|    `OPENID_AUTH_DOMAIN`     |Address of the openid Auth0 domain|
|    `OPENID_CLIENT_ID`     |Auth0 Client ID|
|    `OPENID_CLIENT_SECRET`     |Auth0 Client Secret|
|    `SERVER_COOKIE_DOMAIN`     |Cookie URL domain|
|    `SERVER_REDIRECT_PATH`     |URL to redirect user to after a successful login|
|    `SLACK_BOT_TOKEN`     |Slack token for bot for interactive messages|
|    `SLACK_BOT_ID`     |Slack ID for bot|
|    `SLACK_VERIF_TOKEN`     |Slack OAUTH token for bot for interactive messages|
|    `LDAP_URL`     |URL to connect to LDAP service on|
|    `LDAP_PORT`     |Port to connect to LDAP service on|
|    `LDAP_USER`     |Username to log into LDAP with|
|    `LDAP_PASS`     |Password to log into LDAP with|

The cookies fields are deprecated, but if you want to create a frontend, there are appropriate websocket endpoints. 
You'll need to get the LDAP creds for the program from an RTP.
### Slack integration
If you choose to go the slackbot route, you'll need the bot id, token, and a separate token to verify that the message actually came from slack.
All of that data cna be found at Slack's [bot dashboard](api.slack.com/apps), after you've created your bot.
### Database
Honestly very easy, just set up a PostgreSQL database with username/password that you set variables for. Check out PGSQL docs for help with that.
### Running the app
Make sure you got thos ~~Beans~~ environment variables, then:
1. `go get andreweggleston/DeathByDagger`
2. `cd $(GOPATH)/src/github.com/andreweggleston/DeathByDagger`
## Structure
Everything goes in the folder its named after:

Models in [models](../blob/master/models)
Controllers in [controllers](../blob/master/controllers)
Database in [database](../blob/master/database)
Routes in [routes.go](../blob/master/routes/routes.go)
Helpers in [models](../blob/master/helpers)

## Contributing
1. Fork this project.
2. Make your branch! (`git checkout -b branch_name`)
3. Commit and push your changes. (`git commit` and then `git push origin branch_name`).
4. Make a PR.
Before making a PR, make your code matches Go style guidlines (just run `go fmt`), and squash your commits.
