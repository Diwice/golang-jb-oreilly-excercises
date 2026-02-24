# Testing - Chapter 15
### Task 1
parser (unexported): both num conversion errors, default case; could test them
DataProcessor (exported): basic calculations for all cases; could test each
WriteData (exported): writing from ch to io.Writer; could test if the written value matches
NewController (exported): testing success and bad request cases might be possible; busy case is hardly testable - won't test it
Overall - extra testing with httpserver module

Possible issues: file creation inside main, default case DP totally ignoring absense of operation

Moved all the details not related to server to separate package.
Moved main server to cmd/server.
Removed channel closing in DP, so the channel doesn't close early.
Modified DP to not ignore absense/wrong op.
Modified NewController's Handler to respond with bad request in case of absent/0-length body in order to test bad request case.
### Task 2
No matter how many times I launched the race condition on tests, no matter how many times I sent requests on the built binary with the flag, couldn't catch it at all. 
So, either it remains unfixed in this repo, or I fixed it beforehand in Task 1. 
Can't really tell what exactly had it fixed, could only speculate that closing channel early inside DP caused issues inside WD, since they're separate goroutines and pretty much nothing else there could've caused race condition.
Closing channel early, however, messed up the inputs (request body reads) aswell.. Only removed it because noticed exactly this reason even before launching the code.
Don't know if the author intended it that way, either.
### Task 3
Pretty much no issues here, except the fact that I decided to test DP before parser itself (which the task asks me to do). Well, did that too.
