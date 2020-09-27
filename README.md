# statistic-server
Statistical service - a service that collects indicators using the RESTful API, from other services, aggregates and submits to the requester in the desired form.

The model of collected data consists of a set of indicators:
• Name of the service
• Metric values
• Metric name
• Time
• HTTP status

Implemented the ability to record data and sample data according to certain parameters, using the PostgreSQL database. To test the service, a platform for collaboration and API development - Postman - was used.

In addition to using the basic principles of protection of web-services, the possibility of detecting suspicious HTTP-requests that may indicate a DDoS attack was considered. To model a possible solution to the problem, a single-layer neural network was implemented, which classifies requests for these parameters into suspicious and secure. If the request is suspicious, the bot sends a message in the telegram chat.