package main

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	_ "modernc.org/sqlite"
)

type CommandStatus string

const (
	COMMAND_IDLE    CommandStatus = "IDLE"
	COMMAND_RUNNING CommandStatus = "RUNNING"
	COMMAND_EXITED  CommandStatus = "EXITED"
	COMMAND_ERROR   CommandStatus = "ERROR"
)

type LogFD int

const (
	LOG_STDOUT LogFD = 1
	LOG_STDERR LogFD = 2
	LOG_ERROR  LogFD = -1
)

type CommandInfo struct {
	Id        string
	CreatedAt string
	Command   string
	Status    CommandStatus
}

type CommandLog struct {
	Id        string
	CommandId string
	CreatedAt string
	Content   string
	FD        LogFD
}

type DB interface {
	NewCommand(command string) (*CommandInfo, error)
	GetCommand(id string) (*CommandInfo, error)
	ListCommands(before string, n uint) ([]CommandInfo, error)
	AddLog(log *CommandLog) error
	GetLogs(commandId string, before string, n uint) ([]CommandLog, error)
	UpdateStatus(id string, status CommandStatus) error
}

type CockpitDB struct {
	*sql.DB
}

func NewDB(dataSourceName string) (DB, error) {
	sqlDB, err := sql.Open("sqlite", dataSourceName)
	if err != nil {
		slog.Error("cannot open `cockpit.db` database file", "error", err)
		return nil, err
	}

	db := CockpitDB{
		sqlDB,
	}
	err = db.Init()
	if err != nil {
		return nil, err
	}
	return &db, nil
}

const COMMAND_TABLE_NAME = "command"
const TABLE_SCHEMA_QUERY = "SELECT sql FROM sqlite_schema WHERE name=?"
const CREATE_COMMAND_TABLE_QUERY = `
CREATE TABLE IF NOT EXISTS command (
    id TEXT PRIMARY KEY,
    created_at TEXT NOT NULL,
    command TEXT NOT NULL,
    status TEXT NOT NULL
);
`
const CREATE_LOG_TABLE_QUERY = `
CREATE TABLE IF NOT EXISTS log (
    id TEXT PRIMARY KEY,
    command_id TEXT NOT NULL,
    created_at TEXT NOT NULL,
    content TEXT NOT NULL,
    fd INTEGER NOT NULL,
    FOREIGN KEY (command_id) REFERENCES command (id)
);
`
const INSERT_COMMAND_QUERY = `
INSERT INTO command (id, created_at, command, status)
VALUES (?, ?, ?, ?);
`
const SELECT_COMMAND_QUERY = `
SELECT id, created_at, command, status
FROM command
WHERE id = %1;
`
const LIST_COMMAND_QUERY = `
SELECT id, created_at, command, status 
FROM command
WHERE id < $1
ORDER BY id DESC
LIMIT $2;
`
const UPDATE_STATUS_QUERY = `
UPDATE command
SET status = ?
WHERE id = ?;
`
const INSERT_LOG_QUERY = `
INSERT INTO log (id, command_id, created_at, content, fd)
VALUES (?, ?, ?, ?, ?);
`
const SELECT_LOG_QUERY = `
SELECT id, command_id, created_at, content, fd
FROM log
WHERE command_id = $1 AND id < $2
ORDER BY id DESC
LIMIT $3;
`

func (db *CockpitDB) Init() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		slog.Error("unable to connect to database", "error", err)
		return err
	}

	// create tables if not exist
	if _, err := db.Exec(CREATE_COMMAND_TABLE_QUERY); err != nil {
		slog.Error("unable to create command table", "error", err)
		return err
	}

	if _, err := db.Exec(CREATE_LOG_TABLE_QUERY); err != nil {
		slog.Error("unable to create log table", "error", err)
		return err
	}

	return nil
}

func (db *CockpitDB) NewCommand(command string) (*CommandInfo, error) {
	id := IdGen()
	createdAt := FormatNow()
	status := COMMAND_IDLE
	_, err := db.Exec(INSERT_COMMAND_QUERY, id, createdAt, command, status)
	if err != nil {
		slog.Error("failed to insert new command", "error", err)
		return nil, err
	}

	commandInfo := CommandInfo{
		Id:        id,
		CreatedAt: createdAt,
		Command:   command,
		Status:    status,
	}
	return &commandInfo, nil
}

func (db *CockpitDB) GetCommand(id string) (*CommandInfo, error) {
	var c CommandInfo

	row := db.QueryRow(SELECT_COMMAND_QUERY, id)
	if err := row.Scan(&c.Id, &c.CreatedAt, &c.Command, &c.Status); err != nil {
		return nil, err
	}
	return &c, nil
}

func (db *CockpitDB) ListCommands(before string, n uint) ([]CommandInfo, error) {
	// TODO: think about after and n
	if len(before) == 0 {
		before = MAX_ID
	}

	rows, err := db.Query(LIST_COMMAND_QUERY, before, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var commands []CommandInfo
	for rows.Next() {
		var c CommandInfo
		err := rows.Scan(&c.Id, &c.CreatedAt, &c.Command, &c.Status)
		if err != nil {
			slog.Error("ListCommands", "error", err)
			continue
		}

		commands = append(commands, c)
	}

	if err = rows.Err(); err != nil {
		return commands, err
	}
	return commands, nil
}

func (db *CockpitDB) AddLog(log *CommandLog) error {
	id := log.Id
	commandId := log.CommandId
	createdAt := log.CreatedAt
	content := log.Content
	fd := log.FD

	_, err := db.Exec(INSERT_LOG_QUERY, id, commandId, createdAt, content, fd)
	if err != nil {
		slog.Error("failed to insert new log", "error", err)
		return err
	}

	return nil
}

func (db *CockpitDB) GetLogs(commandId string, before string, n uint) ([]CommandLog, error) {
	if len(before) == 0 {
		before = MAX_ID
	}

	rows, err := db.Query(SELECT_LOG_QUERY, commandId, before, n)
	if err != nil {
		return nil, err
	}

	var logs []CommandLog
	for rows.Next() {
		var l CommandLog
		err := rows.Scan(&l.Id, &l.CommandId, &l.CreatedAt, &l.Content, &l.FD)
		if err != nil {
			slog.Error("GetLogs", "error", err)
			continue
		}

		logs = append(logs, l)
	}

	if err = rows.Err(); err != nil {
		return logs, err
	}
	return logs, nil
}

func (db *CockpitDB) UpdateStatus(id string, status CommandStatus) error {
	_, err := db.Exec(UPDATE_STATUS_QUERY, status, id)
	if err != nil {
		slog.Error("failed to update status", "error", err)
		return err
	}
	return nil
}

