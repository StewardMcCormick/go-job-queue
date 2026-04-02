package main

type App struct{}

func Init() (*App, error) {
	return &App{}, nil
}

func (a *App) Run() error {
	return nil
}

func (a *App) Shutdown() error {
	return nil
}
