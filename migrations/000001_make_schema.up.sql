CREATE TABLE cmd(
    sid uuid DEFAULT gen_random_uuid(),

    status TEXT NOT NULL,

    PRIMARY KEY (sid)
);


CREATE TABLE cmd_info(
    cmd_info_id SERIAL PRIMARY KEY,

    cmd_sid        uuid NOT NULL,

    scheduled_time TIMESTAMP NOT NULL,
    started_time   TIMESTAMP NOT NULL,
    exit_time      TIMESTAMP NOT NULL,

    exit_code      INT NOT NULL,


    CONSTRAINT FK_cmd_info_cmd
        FOREIGN KEY (cmd_sid)
        REFERENCES cmd (sid)
);

CREATE TABLE cmd_source(
    source_id SERIAL PRIMARY KEY,

    cmd_sid     uuid NOT NULL,

    source_line TEXT NOT NULL,

    CONSTRAINT FK_cmd_source_cmd
        FOREIGN KEY (cmd_sid)
        REFERENCES cmd (sid)
);

CREATE TABLE cmd_output(
    output_id SERIAL PRIMARY KEY,

    cmd_sid      uuid NOT NULL,
    output_line  TEXT NOT NULL,

    CONSTRAINT FK_cmd_output_cmd
        FOREIGN KEY (cmd_sid)
        REFERENCES cmd (sid)
);