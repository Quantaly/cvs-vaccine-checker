# cvs-vaccine-checker

For whatever reason, CVS has a JSON API for the availability of their COVID-19 vaccines. Not sure why but also _super_ not complaining, got my appointment because of it. This program, when run periodically, can post to a Discord webhook whenever there are vaccines available in your area.

## Setup

Create a webhook in a Discord channel and copy the URL.

Remove the ".template" extension from [structs/discord-webhook-url.go.template](https://github.com/Quantaly/cvs-vaccine-checker/blob/main/structs/discord-webhook-url.go.template) and paste in your webhook URL for the value of `discordWebhookUrl`.

In [cvs-vaccine-checker.go](https://github.com/Quantaly/cvs-vaccine-checker/blob/main/cvs-vaccine-checker.go), edit `apiEndpoint` to reflect your state and `nearbyCities` to reflect where you're willing to travel.

Set up `cron` or something similar to periodically run the program.

Get vaccinated so we can all eventually go back to normal!!
