package app

import (
	"context"
	
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/widgets/donut"
	"github.com/mum4k/termdash/widgets/segmentdisplay"
	"github.com/mum4k/termdash/widgets/text"
)

// widgets represents the various widgets our ui will display.
type widgets struct {
	donTimer        *donut.Donut
	disType         *segmentdisplay.SegmentDisplay
	txtInfo         *text.Text
	txtTimer        *text.Text
	updateDontTimer chan []int
	updateTxtInfo   chan string
	updateTxtTimer  chan string
	updateTxtType   chan string
}

// update updates the widgets with new data. redrawCh indicates when the app should redraw the screen.
func (w *widgets) update(timer []int, txtType, txtInfo, txtTimer string, redrawCh chan<- bool) {
	if txtInfo != "" {
		w.updateTxtInfo <- txtInfo
	}
	
	if txtType != "" {
		w.updateTxtType <- txtType
	}
	
	if txtTimer != "" {
		w.updateTxtTimer <- txtTimer
	}
	
	if len(timer) > 0 {
		w.updateDontTimer <- timer
	}
	
	redrawCh <- true
}

// newWidget uses helper fucntions to initialize and return a new widget.
func newWidget(ctx context.Context, errorCh chan<- error) (*widgets, error) {
	w := &widgets{
		updateDontTimer: make(chan []int),
		updateTxtType:   make(chan string),
		updateTxtInfo:   make(chan string),
		updateTxtTimer:  make(chan string),
	}
	var err error
	
	w.donTimer, err = newDonut(ctx, w.updateDontTimer, errorCh)
	if err != nil {
		return nil, err
	}
	
	w.disType, err = newSegmentDisplay(ctx, w.updateTxtType, errorCh)
	if err != nil {
		return nil, err
	}
	
	w.txtInfo, err = newText(ctx, w.updateTxtInfo, errorCh)
	if err != nil {
		return nil, err
	}
	
	w.txtTimer, err = newText(ctx, w.updateTxtTimer, errorCh)
	if err != nil {
		return nil, err
	}
	
	return w, nil
}

// newText initializes a new Text widget.
func newText(ctx context.Context, updateText <-chan string, errorCh chan<- error) (*text.Text, error) {
	txt, err := text.New()
	if err != nil {
		return nil, err
	}
	
	// Goroutine to update text.
	go func() {
		for {
			select {
			case t := <-updateText:
				txt.Reset()
				errorCh <- txt.Write(t) // forward any errors to the errorCh if they arise.
			case <-ctx.Done():
				return
			}
		}
	}()
	
	return txt, nil
}

// newText initializes a new Text widget, setting the style.
func newDonut(ctx context.Context, donUpdater <-chan []int, errorCh chan<- error) (*donut.Donut, error) {
	don, err := donut.New(
		donut.Clockwise(),
		donut.CellOpts(cell.FgColor(cell.ColorBlue)),
	)
	
	if err != nil {
		return nil, err
	}
	
	go func() {
		for {
			select {
			case d := <-donUpdater:
				if d[0] <= d[1] {
					errorCh <- don.Absolute(d[0], d[1])
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	
	return don, nil
}

// newSegmentDisplay initializes a new segmentDisplay widget.
func newSegmentDisplay(ctx context.Context, updateText <-chan string, errorCh chan<- error) (*segmentdisplay.SegmentDisplay, error) {
	sd, err := segmentdisplay.New()
	if err != nil {
		return nil, err
	}
	
	go func() {
		for {
			select {
			case t := <-updateText:
				if t == "" {
					t = " "
				}
				errorCh <- sd.Write([]*segmentdisplay.TextChunk{
					segmentdisplay.NewChunk(t),
				})
			case <-ctx.Done():
				return
			}
		}
	}()
	
	return sd, nil
}
