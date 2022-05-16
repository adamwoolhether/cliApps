package app

import (
	"context"
	"image"
	"time"

	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/terminal/tcell"
	"github.com/mum4k/termdash/terminal/terminalapi"

	"github.com/adamwoolhether/cliApps/distributing/pomo/pomodoro"
)

// App is used to instantiate and control the interface. Fields are
// unexported because behavior will be controlled through methods.
type App struct {
	ctx        context.Context
	controller *termdash.Controller
	redrawCh   chan bool
	errorCh    chan error
	term       *tcell.Terminal
	size       image.Point
}

// New instantiates a new App.
func New(config *pomodoro.IntervalConfig) (*App, error) {
	ctx, cancel := context.WithCancel(context.Background())

	quitter := func(k *terminalapi.Keyboard) {
		if k.Key == 'q' || k.Key == 'Q' {
			cancel()
		}
	}

	redrawCh := make(chan bool)
	errorCh := make(chan error)

	// Instantiate widgets and buttons.
	w, err := newWidgets(ctx, errorCh)
	if err != nil {
		return nil, err
	}

	s, err := newSummary(ctx, config, redrawCh, errorCh)
	if err != nil {
		return nil, err
	}

	b, err := newButtonSet(ctx, config, w, s, redrawCh, errorCh)
	if err != nil {
		return nil, err
	}

	// Define a new tcell.Terminal to act as the App's backend.
	term, err := tcell.New()
	if err != nil {
		return nil, err
	}

	// Instantiate a new termdash.Container.
	c, err := newGrid(b, w, s, term)
	if err != nil {
		return nil, err
	}

	// Instantiate a new termdash.Controller.
	controller, err := termdash.NewController(term, c, termdash.KeyboardSubscriber(quitter))
	if err != nil {
		return nil, err
	}

	return &App{
		ctx:        ctx,
		controller: controller,
		redrawCh:   redrawCh,
		errorCh:    errorCh,
		term:       term,
	}, nil
}

// resize will determine if the interface needs to be resized
// and returning early if not.
func (a *App) resize() error {
	if a.size.Eq(a.term.Size()) {
		return nil
	}

	a.size = a.term.Size()
	if err := a.term.Clear(); err != nil {
		return err
	}

	return a.controller.Redraw()
}

// Run is used to run and control the app.
func (a *App) Run() error {
	defer a.term.Close()
	defer a.controller.Close()

	// Define a ticker to check for window resizes.
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	// Take actions based on data arriving from channels.
	for {
		select {
		case <-a.redrawCh:
			if err := a.controller.Redraw(); err != nil {
				return err
			}
		case err := <-a.errorCh:
			if err != nil {
				return err
			}
		case <-a.ctx.Done():
			return nil
		case <-ticker.C:
			if err := a.resize(); err != nil {
				return err
			}
		}
	}
}
