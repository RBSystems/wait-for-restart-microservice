# wait-for-restart-microservice
[![CircleCI](https://img.shields.io/circleci/project/byuoitav/wait-for-restart-microservice.svg)](https://circleci.com/gh/byuoitav/wait-for-restart-microservice) [![Codecov](https://img.shields.io/codecov/c/github/byuoitav/wait-for-restart-microservice.svg)](https://codecov.io/gh/byuoitav/wait-for-restart-microservice) [![Apache 2 License](https://img.shields.io/hexpm/l/plug.svg)](https://raw.githubusercontent.com/byuoitav/wait-for-restart-microservice/master/LICENSE)

[![View in Swagger](http://jessemillar.github.io/view-in-swagger-button/button.svg)](https://byuoitav.github.io/swagger-ui/?url=https://raw.githubusercontent.com/byuoitav/wait-for-restart-microservice/master/swagger.json)

A service designed to monitor the restart status of Crestron devices (by testing a TCP connection) and notify the submitter via a POST when restart is complete.

The service monitors multiple targets concurrently, and will notify the address defined in the submission request when they respond successfully to TCP requests.

The information for the pinging information can be found in the Config file, detailed below.

## Config File
When running the program you can pass in the `-config` flag denoting the location of the JSON config file. If no value is passed in it will default to the provided `./config.json`

```
WaitThreshold is the allowed number of devices being monitored before no wait is imposed between iterations
IterativeTime is the number of seconds to wait between iterations, assuming that the WaitThreshold has not been reached
IndividualTimeout is the number of milliseconds to wait for each individual connection (TCP attempt) before timing out
```
