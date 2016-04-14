# Status Microservice

A service designed to monitor the restart status of devices (by testing a TCP connection) and notify the submitter when they are complete. Designed specifically for Crestron devices.

The service monitors multiple targets concurrently, and will notify the address defined in the submission request when they respond successfully to TCP requests.

## Endpoints
