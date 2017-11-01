wip wip hurray

notes:
HOW DOES A HOST SEND DATA TO RECIEVERS???

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

task 2: passing off work
we can start with work being list of commands to execute
    - ex. [ python scrape.py arg1, python scrape.py arg2]

????
I have to do some work
    - i have to scrape a webpage, parse the html, extract some data, put it in a DB
What is work? How is it structured?
    - multiple calls to a program with different arguments
        - up to the programmer of the work to design it in a way that this is possible
    - i think that works
    - we can handle it similarly to hadoopstream
        - have one program that you are going to run
        - have a list of inputs - we distribute these inputs
            - less robust because we can only have one program
How would we distribute work evenly?
    - How to we know which divided work takes more power than others?
        - Perhaps we can constantly moniter each worker's load somehow?
            - then if it gets too high we can remove some of its future work?
            - maybe we can have each worker check in with the host every few seconds and report its load?
How to put the results back together?
    - Work will be split
    - After sending the result back to the host, what does the host do with it?
        - do something similar to hadoopstream and have it write to output file?
This will be quite complex to test