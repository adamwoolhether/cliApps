package app

import (
	"github.com/mum4k/termdash/align"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/container/grid"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/terminalapi"
)

// newGrid defines a new grid layout with the given buttonSet, widgets, and Terminal.
func newGrid(b *buttonSet, w *widgets, s *summary, t terminalapi.Terminal) (*container.Container, error) {
	builder := grid.Builder{}
	
	// Add first row using percentages to allow dynamic resizing.
	builder.Add(
		grid.RowHeightPerc(30,
			grid.ColWidthPercWithOpts(30,
				[]container.Option{
					container.Border(linestyle.Light),
					container.BorderTitle("Press Q to Quit."),
				},
				// Add inside row in left column.
				grid.RowHeightPerc(80,
					grid.Widget(w.donTimer)),
				grid.RowHeightPercWithOpts(20,
					[]container.Option{
						container.AlignHorizontal(align.HorizontalCenter),
					},
					grid.Widget(w.txtTimer,
						container.AlignHorizontal(align.HorizontalCenter),
						container.AlignVertical(align.VerticalMiddle),
						container.PaddingLeftPercent(49),
					),
				),
			),
			// Add right column with remaining 70%.
			grid.ColWidthPerc(70,
				grid.RowHeightPerc(80,
					grid.Widget(w.disType, container.Border(linestyle.Light)),
				),
				grid.RowHeightPerc(20,
					grid.Widget(w.txtInfo, container.Border(linestyle.Light)),
				),
			),
		),
	)
	
	// Add second row.
	builder.Add(
		grid.RowHeightPerc(10,
			grid.ColWidthPerc(50,
				grid.Widget(b.btStart),
			),
			grid.ColWidthPerc(50,
				grid.Widget(b.btPause),
			),
		),
	)
	
	// Add third row.
	builder.Add(
		grid.RowHeightPerc(60,
			grid.ColWidthPerc(30,
				grid.Widget(s.bcDay,
					container.Border(linestyle.Light),
					container.BorderTitle("Daily Summary (minutes)"),
				),
			),
			grid.ColWidthPerc(70,
				grid.Widget(s.lcWeekly,
					container.Border(linestyle.Light),
					container.BorderTitle("Weekly Summary"),
				),
			),
		),
	)
	
	// Build the layout.
	gridOpts, err := builder.Build()
	if err != nil {
		return nil, err
	}
	
	c, err := container.New(t, gridOpts...)
	if err != nil {
		return nil, err
	}
	
	return c, nil
}
