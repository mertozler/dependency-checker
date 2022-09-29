# Dependency Checker

Dependency Checker checks whether the libraries used in javascript projects are using the current version. Retrieves github repos and mail addresses via a REST api. It sends the queries saved in the Redis database every 24 hours as an e-mail. Specifies out-of-date libraries and sends version information for the current ones.

## Installation

First, run Redis on Docker using the following command, since the project uses the redis database. Make the necessary configurations in the yaml file in the configs in the project. SMTP is used for mail and if you do not enter your e-mail address and password, it will not be able to send mail successfully!

```bash
docker compose up -d
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)
