package handler

import (
	"errors"

	"github.com/m1k1o/neko/server/pkg/types"
	"github.com/m1k1o/neko/server/pkg/types/event"
	"github.com/m1k1o/neko/server/pkg/types/message"
	"github.com/m1k1o/neko/server/pkg/xorg"
)

var (
	ErrIsNotAllowedToHost = errors.New("is not allowed to host")
	ErrIsNotTheHost       = errors.New("is not the host")
	ErrIsAlreadyTheHost   = errors.New("is already the host")
	ErrIsAlreadyHosted    = errors.New("is already hosted")
)

func (h *MessageHandlerCtx) controlRelease(session types.Session) error {
	if !session.Profile().CanHost || session.PrivateModeEnabled() {
		return ErrIsNotAllowedToHost
	}

	if !session.IsHost() {
		return ErrIsNotTheHost
	}

	h.desktop.ResetKeys()
	session.ClearHost()

	return nil
}

func (h *MessageHandlerCtx) controlRequest(session types.Session) error {
	if !session.Profile().CanHost || session.PrivateModeEnabled() {
		return ErrIsNotAllowedToHost
	}

	if session.IsHost() {
		return ErrIsAlreadyTheHost
	}

	if h.sessions.Settings().LockedControls && !session.Profile().IsAdmin {
		return ErrIsNotAllowedToHost
	}

	// if implicit hosting is enabled, set session as host without asking
	if h.sessions.Settings().ImplicitHosting {
		session.SetAsHost()
		return nil
	}

	// if there is no host, set session as host
	host, hasHost := h.sessions.GetHost()
	if !hasHost {
		session.SetAsHost()
		return nil
	}

	// TODO: Some throttling mechanism to prevent spamming.

	// let host know that someone wants to take control
	host.Send(
		event.CONTROL_REQUEST,
		message.SessionID{
			ID: session.ID(),
		})

	return ErrIsAlreadyHosted
}

func (h *MessageHandlerCtx) controlMove(session types.Session, payload *message.ControlPos) error {
	if err := h.controlRequest(session); err != nil && !errors.Is(err, ErrIsAlreadyTheHost) {
		return err
	}

	// handle active cursor movement
	h.desktop.Move(payload.X, payload.Y)
	h.webrtc.SetCursorPosition(payload.X, payload.Y)
	return nil
}

func (h *MessageHandlerCtx) controlScroll(session types.Session, payload *message.ControlScroll) error {
	if err := h.controlRequest(session); err != nil && !errors.Is(err, ErrIsAlreadyTheHost) {
		return err
	}

	// TOOD: remove this once the client is fixed
	if payload.DeltaX == 0 && payload.DeltaY == 0 {
		payload.DeltaX = payload.X
		payload.DeltaY = payload.Y
	}

	h.desktop.Scroll(payload.DeltaX, payload.DeltaY, payload.ControlKey)
	return nil
}

func (h *MessageHandlerCtx) controlButtonPress(session types.Session, payload *message.ControlButton) error {
	if payload.ControlPos != nil {
		if err := h.controlMove(session, payload.ControlPos); err != nil {
			return err
		}
	} else if err := h.controlRequest(session); err != nil && !errors.Is(err, ErrIsAlreadyTheHost) {
		return err
	}

	return h.desktop.ButtonPress(payload.Code)
}

func (h *MessageHandlerCtx) controlButtonDown(session types.Session, payload *message.ControlButton) error {
	if payload.ControlPos != nil {
		if err := h.controlMove(session, payload.ControlPos); err != nil {
			return err
		}
	} else if err := h.controlRequest(session); err != nil && !errors.Is(err, ErrIsAlreadyTheHost) {
		return err
	}

	return h.desktop.ButtonDown(payload.Code)
}

func (h *MessageHandlerCtx) controlButtonUp(session types.Session, payload *message.ControlButton) error {
	if payload.ControlPos != nil {
		if err := h.controlMove(session, payload.ControlPos); err != nil {
			return err
		}
	} else if err := h.controlRequest(session); err != nil && !errors.Is(err, ErrIsAlreadyTheHost) {
		return err
	}

	return h.desktop.ButtonUp(payload.Code)
}

func (h *MessageHandlerCtx) controlKeyPress(session types.Session, payload *message.ControlKey) error {
	if payload.ControlPos != nil {
		if err := h.controlMove(session, payload.ControlPos); err != nil {
			return err
		}
	} else if err := h.controlRequest(session); err != nil && !errors.Is(err, ErrIsAlreadyTheHost) {
		return err
	}

	return h.desktop.KeyPress(payload.Keysym)
}

func (h *MessageHandlerCtx) controlKeyDown(session types.Session, payload *message.ControlKey) error {
	if payload.ControlPos != nil {
		if err := h.controlMove(session, payload.ControlPos); err != nil {
			return err
		}
	} else if err := h.controlRequest(session); err != nil && !errors.Is(err, ErrIsAlreadyTheHost) {
		return err
	}

	return h.desktop.KeyDown(payload.Keysym)
}

func (h *MessageHandlerCtx) controlKeyUp(session types.Session, payload *message.ControlKey) error {
	if payload.ControlPos != nil {
		if err := h.controlMove(session, payload.ControlPos); err != nil {
			return err
		}
	} else if err := h.controlRequest(session); err != nil && !errors.Is(err, ErrIsAlreadyTheHost) {
		return err
	}

	return h.desktop.KeyUp(payload.Keysym)
}

func (h *MessageHandlerCtx) controlTouchBegin(session types.Session, payload *message.ControlTouch) error {
	if err := h.controlRequest(session); err != nil && !errors.Is(err, ErrIsAlreadyTheHost) {
		return err
	}
	return h.desktop.TouchBegin(payload.TouchId, payload.X, payload.Y, payload.Pressure)
}

func (h *MessageHandlerCtx) controlTouchUpdate(session types.Session, payload *message.ControlTouch) error {
	if err := h.controlRequest(session); err != nil && !errors.Is(err, ErrIsAlreadyTheHost) {
		return err
	}
	return h.desktop.TouchUpdate(payload.TouchId, payload.X, payload.Y, payload.Pressure)
}

func (h *MessageHandlerCtx) controlTouchEnd(session types.Session, payload *message.ControlTouch) error {
	if err := h.controlRequest(session); err != nil && !errors.Is(err, ErrIsAlreadyTheHost) {
		return err
	}
	return h.desktop.TouchEnd(payload.TouchId, payload.X, payload.Y, payload.Pressure)
}

func (h *MessageHandlerCtx) controlCut(session types.Session) error {
	if err := h.controlRequest(session); err != nil && !errors.Is(err, ErrIsAlreadyTheHost) {
		return err
	}

	return h.desktop.KeyPress(xorg.XK_Control_L, xorg.XK_x)
}

func (h *MessageHandlerCtx) controlCopy(session types.Session) error {
	if err := h.controlRequest(session); err != nil && !errors.Is(err, ErrIsAlreadyTheHost) {
		return err
	}

	return h.desktop.KeyPress(xorg.XK_Control_L, xorg.XK_c)
}

func (h *MessageHandlerCtx) controlPaste(session types.Session, payload *message.ClipboardData) error {
	if err := h.controlRequest(session); err != nil && !errors.Is(err, ErrIsAlreadyTheHost) {
		return err
	}

	// if there have been set clipboard data, set them first
	if payload != nil && payload.Text != "" {
		if err := h.clipboardSet(session, payload); err != nil {
			return err
		}
	}

	return h.desktop.KeyPress(xorg.XK_Control_L, xorg.XK_v)
}

func (h *MessageHandlerCtx) controlSelectAll(session types.Session) error {
	if err := h.controlRequest(session); err != nil && !errors.Is(err, ErrIsAlreadyTheHost) {
		return err
	}

	return h.desktop.KeyPress(xorg.XK_Control_L, xorg.XK_a)
}
