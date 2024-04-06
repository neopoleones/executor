package postgres

const queryInsertCommandBase = "INSERT INTO cmd(sid, status) VALUES (@outSid, @outStatus)"
const queryInsertCommandSourceLine = "INSERT INTO cmd_source(cmd_sid, source_line) VALUES(@outSid, @outLine)"
const queryInsertCommandOutputLine = "INSERT INTO cmd_output(cmd_sid, output_line) VALUES(@outSid, @outLine)"

const queryInsertCommandInfo = `
	INSERT INTO
		cmd_info(cmd_sid, scheduled_time, started_time, exit_time, exit_code)
	VALUES
	    (@outSid, @schedTime, @startTime, @exitTime, @exitCode);
`

// GetCommands
const queryGetCommandBases = "SELECT sid, status FROM cmd"

const queryGetCommandBaseBySid = "SELECT status FROM cmd WHERE sid=@outSid"

const queryGetCommandInfo = `
	SELECT
		scheduled_time, started_time, exit_time, exit_code
	FROM cmd_info
	WHERE cmd_sid=@outSid
`

const queryGetCommandSourceLines = `
	SELECT
		source_line
	FROM cmd_source
	WHERE cmd_sid=@outSid
	ORDER BY source_id ASC 
`

const queryGetCommandOutputLines = `
	SELECT
		output_line
	FROM cmd_output
	WHERE cmd_sid=@outSid
	ORDER BY output_id ASC 
`

const queryUpdateCommandStatus = `
	UPDATE cmd
		SET status = @outStatus
	WHERE sid=@outSid
`

const queryUpdateCommandInfo = `
	UPDATE cmd_info
		SET scheduled_time = @schedTime,
		    started_time = @startTime,
		    exit_time = @exitTime,
		    exit_code = @exitCode
	WHERE cmd_sid=@outSid
`
