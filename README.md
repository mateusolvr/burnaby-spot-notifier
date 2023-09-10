# burnaby-spot-notifier

This project is intended to find available spots at Recreational Centre Activities in Burbany, BC, and send an email alerting about the availability.

### Prerequisites

-   [GNU Make](https://www.gnu.org/software/make/)
-   [Docker](http://docker.com)
-   [Redis](https://redis.io)

### How to configure

1. Rename the file `config_example.yml` to `config.yml`.
2. Edit the email configurations with username (from), password (pass) and the emails that will receive the alerts (to).
3. Edit Activity with the title of the activity you want to monitor.
4. Edit Redis configuration, replacing the necessary information.

## RUN

```bash
make up
```
