### API
`/status`

    Returns:
        status:  ok
        version: patterned string "vN"
        uptime:  duration

        commands:
            scheduled: number  (increases pre-execute)
            exited:    number  (increases post-execute)

`/cmd/schedule`

    Args:
        script:  string

    Returns:
        sid:        string (ID of scheduled script (aka command)) 
        sched_time: when command was scheduled

`/cmd/list`
 
    Args (all optional):
        sort:    asc/desc
        from:    datetime
        before:  datetime

    Returns:
        [] runnable:
            sid:        string (uid?)
            script:     string
            status:     scheduled/done            

            info:
                sched_time: datetime
                exit_time:  datetime
                start_Time: datetime

                code:       number
                output:     []string
                        

`/cmd/get`
    
    Args (URL Param):
        sid: string (uid?)

    Returns:
        found: bool
        
        runnable (optional):
            sid:        string (uid?)
            source:     string
            status:     scheduled/done            

            info:
                sched_time: datetime
                exit_time:  datetime
                start_Time: datetime

                code:       number
                output:     []string


slog?


import (
"github.com/sirupsen/logrus"
"os"
"os/exec"
)

func Execute(script string, command []string) (bool, error) {

    cmd := &exec.Cmd{
        Path:   script,
        Args:   command,
        Stdout: os.Stdout,
        Stderr: os.Stderr,
    }

    c.logger.Info("Executing command ", cmd)

    err := cmd.Start()
    if err != nil {
        return false, err
    }

    err = cmd.Wait()
    if err != nil {
        return false, err
    }

    return true, nil
}

command := []string{
"/<path>/yourscript.sh",
"arg1=val1",
"arg2=val2",
}

Execute("/<path>/yourscript.sh", command)