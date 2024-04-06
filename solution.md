### API
`/status`

    Returns:
        status:  ok
        version: patterned string "vN"
        uptime:  duration

        commands:
            running:   number  (increases pre-execute)
            exited:    number  (increases post-execute)

`/cmd/schedule`

    Args(json):
        script:  []string

    Returns:
        status:  ok
        sid:        string (ID of scheduled script (aka command)) 
        sched_time: when command was scheduled

`/cmd/list`
    
    Returns:
        status: ok
        commands: [] runnable:
            sid:        string (uid?)
            status:     scheduled/done
                        

`/cmd/get`
    
    Args (URL Param):
        sid: string (uid?)

    Returns:
        status: ok
        
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