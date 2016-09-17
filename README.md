# go-pool-manager
- This project uses the Radix.v2 library for Redis connection available at https://github.com/mediocregopher/radix.v2
as recommended by Redis itself http://redis.io/clients#go

- it makes use of Pool object which is the thread safe connection pool of the project.
- it works by launching total_processors go rountines into the wild and let them do their own logic in retrieving an event ID and assign it to the shell script. if the id returned is 0 which is converted to -1 it means the list doesn't contain any IDs and the job is done.

- the code will look for  a Redis list names testlist assuming it contains the list of IDs it needs.
