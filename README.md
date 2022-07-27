# Bidoof
![Bidoof in water, holding a stick in its mouth](/docs/bidoof-chew.jpeg "a bidoof")

This is a just-for-fun app to learn some more Go(lang).

It just watches a message queue and does a few things.

I wanted it to be Pokemon themed, and I think the noble Bidoof is the closest to a gopher Pokemon there is.
Also, it does do some beaver like activities (builds dams, swims, etc) so I'm guessing it can also "chew" through the message queue 😉

For this to work properly you'll need
* Amqp message queue (I used RabbitMQ)
* Mysql DB with an "attendees" table (see the sql file in the `docs` directory)
* Copy the `.env_example` file into `.env` and fill out the details (the queue name you specify should already exist)
* run `go run mqchew.go`

As the messages roll in (and assuming they are formatted properly) the app will log out what it's doing.
![Screenshot of application processing a message](/docs/helpful-screenshot.png "a screenshot")
