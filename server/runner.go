package main

type Runner interface {
	NewCommand(db *DB) error
	Run(db *DB, command CommandInfo) error
}

func NewRunner() Runner {
	runner := CockpitRunner{}
	return &runner
}

type CockpitRunner struct{}

func (r *CockpitRunner) NewCommand(db *DB) error {
	return nil
}

func (r CockpitRunner) Run(db *DB, command CommandInfo) error {
	return nil
}
