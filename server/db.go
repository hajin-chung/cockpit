package main

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

type CommandInfo struct {
	Id        string
	CreatedAt string
	Command   string
	Status    string
}

type CommandLog struct {
	CreatedAt string
	Content   string
	FD        uint
}

type DB interface {
	NewCommand(command string) error
	ListCommands(after string, n uint) ([]CommandInfo, error)
	AddLog(id string, log string) error
	UpdateStatus(id string, status string) error
}

func NewDB() DB {
	db := CockpitDB{}
	return &db
}

type CockpitDB struct {
	conn *sql.DB
}

func (d *CockpitDB) NewCommand(command string) error {
	return nil
}

func (d *CockpitDB) ListCommands(after string, n uint) ([]CommandInfo, error) {
	return []CommandInfo{}, nil
}

func (d *CockpitDB) AddLog(id string, log string) error {
	return nil
}

func (d *CockpitDB) UpdateStatus(id string, status string) error {
	return nil
}
