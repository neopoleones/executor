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



chi as lightweight

system exec environment from stdlib

custom buffer for defeating RC 

yaml used for configuration 
    (also as lightweight solution) 