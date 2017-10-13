wip wip hurray

notes:

goal:
Have a "master" node that connects with the outside world - host
Spawn some children nodes (can be hosts/workers) that the master can connect to via tcp sockets
    - make some sort of protocol? - will prolly have to to distribute and do the work
Have the master distribute (?) work between the children nodes
    - if the child is a host, have it distribute the work further, to its children
Worker nodes do work, compute some "result" then send it back up the chain of parents until it hits master
    - master waits (?) on all work to be completed, then sends back to requester via tcp (?)

task 1: spawning instances
multiple workers at once? allow to specify port then pattern to add hosts on? like port 8000, add1
how would do for hosts? - idk cause we need everything to connect back to each other... zzzz