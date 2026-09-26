Distributed Job Queue : 

Multiple Workers -> Queue [job1,job2,job3] -> Processes Job

At high level, Workers atomically dequeue jobs from queue and process them and ack
and then move to completed job 

Job can be scheduled->ready->processing->completed->failed->terminated

At each point,we need ready jobs by time interval and need to fetch the newly scheduled jobs and requeue jobs periodically on fail and frequent removal,sorting,updating,deleting

priority_queue + indexing and having 3 queues determine state : 
ready_queue,in_flight job,delayed_queue

We include custom max_attemps per Job,as each job defines its own attempt,Config
Exponential backoffs and reattempt and if exceeded -> DLQ

Workers,Queues can crash/go offline

We need way to know worker crashed,etc.
Hearbeat mechanism

Need to know whether worker should continue or not :
lease mechanism

Decision :
We give a finite time interval within with the worker if doesnt respond,to drop off the 
current job , and requeue for another worker (includes save point state)

need reliability and consistency with atleast once delivery per job
preventing duplicate job processing (idempotent)

architectural decision : 
include the save_point progress into the Job as defined by the external service
everytime each worker completes job by moving and saving through save_points as defined 
and completes job 
iF fails at any point,another worker comes and continues from the save point and 
any older savepoints (single job transaction ) are prevents by fencing tokens 

job.go : 
includes job states,and the worker related info(if any worker processing)

clock.go :
general clock interface - later using fakeclock can simulate testing

indexed_heap.go :
indexed heap => priority_queue with map generic impl for the 3types of queues

broker.go :
defines all the queues,dead jobs,observability(later) and other func requirements


