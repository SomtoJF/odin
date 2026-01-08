package components

import (
	"fmt"

	"github.com/SomtoJF/odin/core/api"
	"github.com/SomtoJF/odin/core/state"
	introascii "github.com/SomtoJF/odin/ui/components/introascii"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// UI represents the terminal user interface
type UI struct {
	app *tview.Application

	// Backend communication channels (injected)
	requestChan  chan<- api.Request       // Send requests to backend
	responseChan <-chan api.StateSnapshot // Receive state snapshots from backend
	updateChan   <-chan api.StateUpdate   // Receive real-time updates from backend

	// UI components
	messageView *tview.TextView
	inputField  *tview.InputField
	todoList    *tview.List
	statusBar   *tview.TextView
	layout      *tview.Flex
	banner      *tview.TextView
	// UI state
	currentMode state.AgentMode
}

// NewUI creates a new UI instance with injected channels
func NewUI(
	requestChan chan<- api.Request,
	responseChan <-chan api.StateSnapshot,
	updateChan <-chan api.StateUpdate,
) *UI {
	return &UI{
		app:          tview.NewApplication(),
		requestChan:  requestChan,
		responseChan: responseChan,
		updateChan:   updateChan,
		currentMode:  state.AgentModeAsk, // Default mode
	}
}

// Run starts the UI and blocks until the application exits
func (ui *UI) Run() error {
	ui.initializeComponents()
	ui.buildLayout()
	ui.setupEventHandlers()

	// Start listening to backend updates
	go ui.listenToBackend()

	// Run the application (blocks)
	return ui.app.SetRoot(ui.layout, true).Run()
}

// initializeComponents creates all tview components
func (ui *UI) initializeComponents() {
	// Intro ASCII banner
	introComponent := introascii.NewComponent()
	ui.banner = tview.NewTextView().
		SetText(introComponent.IntroASCII()).
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(false)
	// banner.SetBorder(true).SetTitle("Odin Code")

	// Message view (conversation history)
	ui.messageView = tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetChangedFunc(func() {
			ui.app.Draw()
		})
	ui.messageView.SetBorder(true).SetTitle("Messages")

	// Todo list
	ui.todoList = tview.NewList().
		ShowSecondaryText(false)
	ui.todoList.SetBorder(true).SetTitle("TODOs")

	// Input field
	ui.inputField = tview.NewInputField().
		SetLabel("> ").
		SetFieldWidth(0)
	ui.inputField.SetBorder(true).SetTitle("Input")

	// Status bar
	ui.statusBar = tview.NewTextView().
		SetDynamicColors(true).
		SetText(fmt.Sprintf("[Mode: %s] [Status: Idle]", ui.currentMode))
	ui.statusBar.SetBorder(true).SetTitle("Status")
}

// buildLayout creates the UI layout structure
func (ui *UI) buildLayout() {
	// Right panel (todos + status)
	rightPanel := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(ui.todoList, 0, 1, false).
		AddItem(ui.statusBar, 3, 0, false)

	// Main content (messages + input)
	mainPanel := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(ui.messageView, 0, 1, false).
		AddItem(ui.inputField, 3, 0, true)

	mainSection := tview.NewFlex().
		AddItem(mainPanel, 0, 3, true).
		AddItem(rightPanel, 40, 0, false)

	// Overall layout
	ui.layout = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(ui.banner, 0, 1, false).
		AddItem(mainSection, 0, 1, false)
}

// setupEventHandlers sets up keyboard and input handlers
func (ui *UI) setupEventHandlers() {
	// Handle Enter key in input field
	ui.inputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			message := ui.inputField.GetText()
			if message != "" {
				ui.handleUserInput(message)
				ui.inputField.SetText("")
			}
		}
	})

	// Global key handlers for mode switching
	ui.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'a':
			if event.Modifiers()&tcell.ModCtrl != 0 {
				ui.currentMode = state.AgentModeAsk
				ui.updateStatusBar()
				return nil
			}
		case 'e':
			if event.Modifiers()&tcell.ModCtrl != 0 {
				ui.currentMode = state.AgentModeEdit
				ui.updateStatusBar()
				return nil
			}
		case 'p':
			if event.Modifiers()&tcell.ModCtrl != 0 {
				ui.currentMode = state.AgentModePlan
				ui.updateStatusBar()
				return nil
			}
		case 'q':
			if event.Modifiers()&tcell.ModCtrl != 0 {
				ui.app.Stop()
				return nil
			}
		}
		return event
	})
}

// handleUserInput sends user input to the backend
func (ui *UI) handleUserInput(message string) {
	// Send request to backend via channel
	ui.requestChan <- api.Request{
		Message: message,
		Mode:    ui.currentMode,
	}
}

// listenToBackend listens for updates from the backend
func (ui *UI) listenToBackend() {
	for {
		select {
		case snapshot := <-ui.responseChan:
			ui.updateUIFromSnapshot(snapshot)
		case update := <-ui.updateChan:
			ui.handleUpdate(update)
		}
	}
}

// updateUIFromSnapshot updates all UI components from a state snapshot
func (ui *UI) updateUIFromSnapshot(snapshot api.StateSnapshot) {
	ui.app.QueueUpdateDraw(func() {
		ui.updateMessageView(snapshot.Messages)
		ui.updateTodoList(snapshot.Messages)
		ui.updateStatusBar()
	})
}

// updateMessageView updates the message view with conversation history
func (ui *UI) updateMessageView(messages []state.Message) {
	messageText := ""
	for _, msg := range messages {
		messageText += fmt.Sprintf("[yellow]User:[white] %s\n", msg.Body)
		if msg.AnswerSummary != "" {
			messageText += fmt.Sprintf("[green]Agent:[white] %s\n\n", msg.AnswerSummary)
		}

		// Show real-time updates if any
		for _, update := range msg.Updates {
			messageText += fmt.Sprintf("[cyan]... %s\n", update)
		}
	}
	ui.messageView.SetText(messageText)
	ui.messageView.ScrollToEnd()
}

// updateTodoList updates the todo list from the latest message
func (ui *UI) updateTodoList(messages []state.Message) {
	ui.todoList.Clear()

	if len(messages) == 0 {
		return
	}

	// Show todos from the latest message
	lastMsg := messages[len(messages)-1]
	for _, todo := range lastMsg.Todos {
		status := "[ ]"
		color := "white"

		switch todo.Status {
		case state.TodoStatusInProgress:
			status = "[~]"
			color = "yellow"
		case state.TodoStatusCompleted:
			status = "[✓]"
			color = "green"
		}

		ui.todoList.AddItem(
			fmt.Sprintf("[%s]%s %s", color, status, todo.Content),
			"",
			0,
			nil,
		)
	}
}

// updateStatusBar updates the status bar with current mode and execution state
func (ui *UI) updateStatusBar() {
	// Request current snapshot to get execution state
	// (This will be updated by the next snapshot from the backend)
	statusText := fmt.Sprintf(
		"[Mode: [cyan]%s[white]] | Ctrl+A: Ask | Ctrl+E: Edit | Ctrl+P: Plan | Ctrl+Q: Quit",
		ui.currentMode,
	)
	ui.statusBar.SetText(statusText)
}

// handleUpdate handles real-time updates from the backend
func (ui *UI) handleUpdate(update api.StateUpdate) {
	ui.app.QueueUpdateDraw(func() {
		switch update.Type {
		case api.UpdateTypeMessage:
			// Message update handled by snapshot
		case api.UpdateTypeTodo:
			// Todo update handled by snapshot
		case api.UpdateTypeStatus:
			// Status update handled by snapshot
		case api.UpdateTypeToolHistory:
			// Could display tool execution in a separate panel
		}
	})
}
